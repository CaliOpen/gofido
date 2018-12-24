package store

import (
    "github.com/CaliOpen/gofido/config"
    log "github.com/Sirupsen/logrus"
    "github.com/gocql/gocql"
    "github.com/scylladb/gocqlx"
    "github.com/scylladb/gocqlx/qb"
    "github.com/scylladb/gocqlx/table"
    "github.com/tstranex/u2f"
    "time"
)

// Cassandra challenge model
type Challenge struct {
    User          string
    Challenge     string
    Timestamp     time.Time `db:"expire_at"`
    AppId         string
    TrustedFacets []string
}

var challengeMetadata = table.Metadata{
    Name:    "challenge",
    Columns: []string{"user", "challenge", "expire_at", "app_id", "trusted_facets"},
    PartKey: []string{"user", "challenge"},
}

// Cassandra registration model
type Registration struct {
    User      string
    KeyHandle []byte `db:"key_handle"`
    // Raw serialized registration data as received from the token.
    Data []byte
}

var registrationMetadata = table.Metadata{
    Name:    "registration",
    Columns: []string{"user", "key_handle", "data"},
    PartKey: []string{"user", "key_handle"},
}

// Store challenge emmited for a key for validation
type KeyChallenge struct {
    User      string
    KeyHandle []byte `db:"key_handle"`
    Challenge []byte
}

var keyChallengeMetadata = table.Metadata{
    Name:    "key_challenge",
    Columns: []string{"user", "key_handle", "challenge"},
    PartKey: []string{"user", "key_handle"},
}

// Store key counter value
type KeyCounter struct {
    User      string
    KeyHandle []byte `db:"key_handle"`
    Counter   uint32
}

var keyCounterMetadata = table.Metadata{
    Name:    "key_counter",
    Columns: []string{"user", "key_handle", "counter"},
    PartKey: []string{"user", "key_handle"},
}

// Reflect structure and table for easy CRUD
var challengeTable = table.New(challengeMetadata)
var registrationTable = table.New(registrationMetadata)
var keyChallengeTable = table.New(keyChallengeMetadata)
var keyCounterTable = table.New(keyCounterMetadata)

// Our struct for store management
type FidoStore struct {
    Session      *gocql.Session
    ChallengeTtl time.Duration
    AppId        string
}

// Initialize the store engine
func (store *FidoStore) Initialize(config *config.Config) error {
    cluster := gocql.NewCluster(config.Store.Hosts...)
    cluster.Keyspace = config.Store.Keyspace
    cluster.Consistency = config.Store.Consistency
    session, err := cluster.CreateSession()
    if err != nil {
        log.Error("Unable to create sesssion to cassandra store ", err)
        return err
    }
    store.Session = session
    ttl, err := time.ParseDuration(config.Server.ChallengeTtl)
    if err != nil {
        log.Error("Invalid duration value for challenge_ttl ", config.Server.ChallengeTtl)
        return err
    }
    store.ChallengeTtl = ttl
    store.AppId = config.Server.AppId
    log.Info("Store initialized with keyspace ", cluster.Keyspace)
    return nil
}

// Create a new challenge, store it and return an u2f.Challenge
func (store *FidoStore) NewChallenge(user string) (u2f.Challenge, error) {
    challenge := &Challenge{User: user}
    trustedFacets := []string{store.AppId}

    u2f_challenge, err := u2f.NewChallenge(store.AppId, trustedFacets)
    if err != nil {
        return u2f.Challenge{}, err
    }
    challenge.Challenge = EncodeBase64(u2f_challenge.Challenge)
    // TOFIX make challenge validity configurable
    challenge.Timestamp = time.Now().Add(store.ChallengeTtl)
    challenge.AppId = store.AppId
    challenge.TrustedFacets = trustedFacets

    builder := qb.Insert(challengeTable.Name()).Columns(challengeMetadata.Columns...).TTL(store.ChallengeTtl)
    stmt, names := builder.ToCql()

    q := gocqlx.Query(store.Session.Query(stmt), names).BindStruct(challenge)
    if err := q.ExecRelease(); err != nil {
        return u2f.Challenge{}, err
    }
    return *u2f_challenge, err
}

func (store *FidoStore) marshalChallenge(source Challenge) u2f.Challenge {
    challenge_id, _ := DecodeBase64(source.Challenge)
    dest := u2f.Challenge{Challenge: challenge_id,
        AppID:         source.AppId,
        TrustedFacets: source.TrustedFacets,
        Timestamp:     source.Timestamp,
    }
    return dest
}

// Get a challenge by its identifier
func (store *FidoStore) GetChallenge(user string, challenge_id string) (u2f.Challenge, error) {
    challenge := &Challenge{}
    builder := qb.Select(challengeTable.Name()).Columns(challengeMetadata.Columns...)
    builder = builder.Where(qb.Eq("user"), qb.Eq("challenge"))
    stmt, names := builder.ToCql()

    q := gocqlx.Query(store.Session.Query(stmt, user, challenge_id), names)
    if err := q.GetRelease(challenge); err != nil {
        return u2f.Challenge{}, err
    }
    return store.marshalChallenge(*challenge), nil
}

func (store *FidoStore) NewRegistration(user string, challenge u2f.Challenge, resp u2f.RegisterResponse) (*u2f.Registration, error) {

    config := &u2f.Config{
        // Chrome 66+ doesn't return the device's attestation
        // certificate by default.
        SkipAttestationVerify: true,
    }

    u2f_reg, err := u2f.Register(resp, challenge, config)
    if err != nil {
        return &u2f.Registration{}, err
    }

    data, err := u2f_reg.MarshalBinary()
    if err != nil {
        log.Error("Marshall registration failed: ", err)
        return &u2f.Registration{}, err
    }

    // Map u2f registration to store one and insert it
    reg := &Registration{User: user}
    reg.KeyHandle = u2f_reg.KeyHandle
    reg.Data = data

    stmt, names := registrationTable.Insert()
    q := gocqlx.Query(store.Session.Query(stmt), names).BindStruct(reg)
    if err := q.ExecRelease(); err != nil {
        return &u2f.Registration{}, err
    }
    if err = store.insertKeyCounter(user, u2f_reg.KeyHandle); err != nil {
        return &u2f.Registration{}, err
    }
    return u2f_reg, nil
}

func (store *FidoStore) marshallRegistration(source *Registration) u2f.Registration {
    var u2f_reg u2f.Registration
    err := u2f_reg.UnmarshalBinary(source.Data)
    if err != nil {
        log.Error("Unmarshal regitration failed: ", err)
        return u2f.Registration{}
    }
    return u2f_reg
}

func (store *FidoStore) GetRegistrations(username string) ([]u2f.Registration, error) {
    var regs []*Registration
    builder := qb.Select(registrationTable.Name()).Columns(registrationMetadata.Columns...).Where(qb.Eq("user"))
    stmt, names := builder.ToCql()

    q := gocqlx.Query(store.Session.Query(stmt, username), names)
    if err := q.SelectRelease(&regs); err != nil {
        return []u2f.Registration{}, err
    }
    registrations := make([]u2f.Registration, len(regs))
    for i, reg := range regs {
        registrations[i] = store.marshallRegistration(reg)
    }
    return registrations, nil
}

func (store *FidoStore) InsertKeyChallenge(user string, key_handle []byte, challenge u2f.Challenge) error {
    key_challenge := &KeyChallenge{User: user, KeyHandle: key_handle, Challenge: challenge.Challenge}
    builder := qb.Insert(keyChallengeTable.Name()).Columns(keyChallengeMetadata.Columns...).TTL(store.ChallengeTtl)
    stmt, names := builder.ToCql()

    q := gocqlx.Query(store.Session.Query(stmt), names).BindStruct(key_challenge)
    if err := q.ExecRelease(); err != nil {
        return err
    }
    return nil
}

func (store *FidoStore) GetKeyChallenges(user string, key_handle []byte) []KeyChallenge {
    var key_challenges []KeyChallenge
    builder := qb.Select(keyChallengeTable.Name()).Columns(keyChallengeMetadata.Columns...).Where(qb.Eq("user"), qb.Eq("key_handle"))
    stmt, names := builder.ToCql()
    q := gocqlx.Query(store.Session.Query(stmt, user, key_handle), names)
    if err := q.SelectRelease(&key_challenges); err != nil {
        return []KeyChallenge{}
    }
    return key_challenges
}

// Create a new counter for key
func (store *FidoStore) insertKeyCounter(user string, key_handle []byte) error {
    counter := &KeyCounter{User: user, KeyHandle: key_handle, Counter: 0}
    stmt, names := keyCounterTable.Insert()
    q := gocqlx.Query(store.Session.Query(stmt), names).BindStruct(counter)
    if err := q.ExecRelease(); err != nil {
        return err
    }
    return nil
}

func (store *FidoStore) GetKeyCounter(user string, key_handle []byte) (KeyCounter, error) {
    var counter KeyCounter
    builder := qb.Select(keyCounterTable.Name()).Columns(keyCounterMetadata.Columns...)
    builder = builder.Where(qb.Eq("user"), qb.Eq("key_handle"))
    stmt, names := builder.ToCql()
    q := gocqlx.Query(store.Session.Query(stmt, user, key_handle), names)
    if err := q.GetRelease(&counter); err != nil {
        return KeyCounter{}, err
    }
    return counter, nil
}

func (store *FidoStore) UpdateCounter(user string, key_handle []byte, counter uint32) error {
    builder := qb.Update(keyCounterTable.Name()).Set("counter").Where(qb.Eq("user"), qb.Eq("key_handle"))
    stmt, names := builder.ToCql()
    q := gocqlx.Query(store.Session.Query(stmt, counter, user, key_handle), names)
    if err := q.ExecRelease(); err != nil {
        return err
    }
    return nil
}
