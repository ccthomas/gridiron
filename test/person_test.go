package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/ccthomas/gridiron/internal/person"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewPerson(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")

	unique := uuid.New().String()
	p := person.CreateNewPersonDTO{
		Name: fmt.Sprintf("TestTeamName%s", unique),
	}

	// When

	res, actual := sendApiReq[person.Person](
		t,
		http.MethodPost,
		"http://localhost:8080/person",
		p,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.NotNil(t, actual.Id, "Actual id is not nil.")
	assert.Equal(t, p.Name, actual.Name, "Actual name is incorrect.")
	assert.Equal(t, tn.Id, actual.TenantId, "Actual tenant id is incorrect.")
}

func TestNewPersonContract_TypeAthlete(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")
	tm := createTeam(t, tn.Id, "TestTeamName")
	p := createPerson(t, tn.Id, "TestPersonName")

	pc := person.CreateNewPersonContractDTO{
		PersonId:   p.Id,
		EntityId:   tm.Id,
		EntityType: person.Team,
		Type:       person.Athlete,
	}

	// When

	res, actual := sendApiReq[person.PersonContract](
		t,
		http.MethodPost,
		"http://localhost:8080/person/contract",
		pc,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.NotNil(t, actual.Id, "Actual id is not nil.")
	assert.Equal(t, p.Id, actual.PersonId, "Actual id is incorrect.")
	assert.Equal(t, tn.Id, actual.TenantId, "Actual tenant id is incorrect.")
	assert.Equal(t, tm.Id, actual.EntityId, "Actual id is incorrect.")
	assert.Equal(t, pc.EntityType, actual.EntityType, "Actual entity type is incorrect.")
	assert.Equal(t, pc.Type, actual.Type, "Actual type is incorrect.")
}

func TestNewPersonContract_TypeCoach(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")
	tm := createTeam(t, tn.Id, "TestTeamName")
	p := createPerson(t, tn.Id, "TestPersonName")

	pc := person.CreateNewPersonContractDTO{
		PersonId:   p.Id,
		EntityId:   tm.Id,
		EntityType: person.Team,
		Type:       person.Coach,
	}

	// When

	res, actual := sendApiReq[person.PersonContract](
		t,
		http.MethodPost,
		"http://localhost:8080/person/contract",
		pc,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.NotNil(t, actual.Id, "Actual id is not nil.")
	assert.Equal(t, p.Id, actual.PersonId, "Actual id is incorrect.")
	assert.Equal(t, tn.Id, actual.TenantId, "Actual tenant id is incorrect.")
	assert.Equal(t, tm.Id, actual.EntityId, "Actual id is incorrect.")
	assert.Equal(t, pc.EntityType, actual.EntityType, "Actual entity type is incorrect.")
	assert.Equal(t, pc.Type, actual.Type, "Actual type is incorrect.")
}

func TestNewPersonContract_TypeOwner(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")
	tm := createTeam(t, tn.Id, "TestTeamName")
	p := createPerson(t, tn.Id, "TestPersonName")

	pc := person.CreateNewPersonContractDTO{
		PersonId:   p.Id,
		EntityId:   tm.Id,
		EntityType: person.Team,
		Type:       person.Owner,
	}

	// When

	res, actual := sendApiReq[person.PersonContract](
		t,
		http.MethodPost,
		"http://localhost:8080/person/contract",
		pc,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.NotNil(t, actual.Id, "Actual id is not nil.")
	assert.Equal(t, p.Id, actual.PersonId, "Actual id is incorrect.")
	assert.Equal(t, tn.Id, actual.TenantId, "Actual tenant id is incorrect.")
	assert.Equal(t, tm.Id, actual.EntityId, "Actual id is incorrect.")
	assert.Equal(t, pc.EntityType, actual.EntityType, "Actual entity type is incorrect.")
	assert.Equal(t, pc.Type, actual.Type, "Actual type is incorrect.")
}

func TestGetPersonByTeamContract_EmptyData(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")

	// When

	res, actual := sendApiReq[person.PersonGetDTO](
		t,
		http.MethodPost,
		fmt.Sprintf("http://localhost:8080/person/contract/team/%s", uuid.New().String()),
		nil,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.Equal(t, 0, actual.Count, "count does not equal 0.")
	assert.Equal(t, 0, len(actual.Data), "data count does not equal 0.")
}

func TestGetPersonByTeamContract(t *testing.T) {
	// Given

	u, loginRes := login(t)
	tn := createTenant(t, u.Id, "TestTenantName")
	tm := createTeam(t, tn.Id, "TestTeamName")
	p1 := createPerson(t, tn.Id, "TestPersonName")
	p2 := createPerson(t, tn.Id, "TestPersonName")

	pc1 := createPersonContract(t, tn.Id, person.CreateNewPersonContractDTO{
		PersonId:   p1.Id,
		EntityId:   tm.Id,
		EntityType: person.Team,
		Type:       person.Owner,
	})

	pc2 := createPersonContract(t, tn.Id, person.CreateNewPersonContractDTO{
		PersonId:   p2.Id,
		EntityId:   tm.Id,
		EntityType: person.Team,
		Type:       person.Owner,
	})

	// When

	res, actual := sendApiReq[person.PersonGetDTO](
		t,
		http.MethodPost,
		fmt.Sprintf("http://localhost:8080/person/contract/team/%s", tm.Id),
		nil,
		loginRes.AccessToken,
		tn.Id,
	)

	// Then

	assert.Equal(t, http.StatusOK, res.StatusCode, "Status code is not a 200")
	assert.Equal(t, 2, actual.Count, "count does not equal 0.")
	assert.Equal(t, 2, len(actual.Data), "data count does not equal 0.")
	assert.Equal(t, person.PersonWithContractDTO{
		Person:   p1,
		Contract: pc1,
	}, actual.Data[0], "First data is not correct.")
	assert.Equal(t, person.PersonWithContractDTO{
		Person:   p2,
		Contract: pc2,
	}, actual.Data[1], "Second data is not correct.")
}
