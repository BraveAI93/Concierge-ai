package intelligence_runtime

import "context"

// IdentityProvisioner is the high-trust boundary for initial durable identity
// creation. Runtime ingestion and the kernel repository do not implement this
// port. Rebinding or deleting a stable subject/profile is intentionally absent.
type IdentityProvisioner interface {
	Provision(ctx context.Context, binding PersonBinding) error
}

// PostgresIdentityProvisioner uses the separate ci_kernel_identity_provisioner
// role. Local staging grants membership to the migration executor solely to
// test the boundary; production must use a separately authorized backend
// principal and must not grant this role to ordinary runtime credentials.
type PostgresIdentityProvisioner struct {
	Repository *PostgresRuntimeRepository
}

func (p PostgresIdentityProvisioner) Provision(ctx context.Context, binding PersonBinding) error {
	if p.Repository == nil || binding.Person.ID == "" || binding.World.ID == "" || binding.World.PersonID != binding.Person.ID || binding.StableSubjectID == "" || binding.SourceProfileID == "" {
		return ErrInvalidRuntimeConfig
	}
	tx, err := p.Repository.beginAsRole(ctx, identityProvisionerRole, "", binding.StableSubjectID)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	personPayload, err := jsonPayload(binding.Person)
	if err != nil {
		return err
	}
	worldPayload, err := jsonPayload(binding.World)
	if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.people(id,payload) VALUES($1,$2::jsonb)", binding.Person.ID, string(personPayload)); isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	} else if err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.worlds(person_id,payload) VALUES($1,$2::jsonb)", binding.Person.ID, string(worldPayload)); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.person_binding_subjects(stable_subject,person_id) VALUES($1,$2)", binding.StableSubjectID, binding.Person.ID); isUniqueViolation(err) {
		return ErrDuplicateRuntimeRecord
	} else if err != nil {
		return err
	}
	profiles := append([]string{binding.SourceProfileID}, binding.AllowedSourceProfileIDs...)
	seen := make(map[string]bool)
	for index, profile := range profiles {
		if profile == "" || seen[profile] {
			continue
		}
		seen[profile] = true
		if _, err = tx.ExecContext(ctx, "INSERT INTO ci_kernel_v04.person_profile_links(person_id,source_profile_id,is_primary) VALUES($1,$2,$3)", binding.Person.ID, profile, index == 0); isUniqueViolation(err) {
			return ErrDuplicateRuntimeRecord
		} else if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Bindings are immutable initial identity mappings. The explicit runtime
// repository method is retained only for source compatibility and fails closed;
// callers must use a separately authorized IdentityProvisioner.
func (r *PostgresRuntimeRepository) SeedBinding(context.Context, PersonBinding) error {
	return ErrIdentityProvisioningRequired
}

// Compile-time proof that provisioning is outside the ordinary runtime ports.
var _ IdentityProvisioner = PostgresIdentityProvisioner{}
