# Changelog

## v1.0.2

### Patch Release

### Fixed

- Corrected route registration issues.
- Fixed authorization checks for job ownership.
- Restricted customer/staff operations according to RBAC.
- Improved regression-tested authorization behavior.
- Minor routing and validation fixes.

### Notes

- End-to-end regression testing completed successfully for v1.0.2.
- Fresh database installation verified.



## v1.1.0

### Added

* Added endpoint to retrieve a customer by ID.
* Added endpoint to retrieve a staff member by ID.
* Added endpoint to retrieve a job by its application-generated Job Number.

### Changed

* Job retrieval now uses the application-generated Job Number as the public identifier instead of the MongoDB ObjectID.
* Restored intended authorization allowing Receptionists, Technicians, Head Technicians, Administrators, and Super Administrators to view any job.
* Customers remain restricted to viewing only their own jobs.
* Job number lookups are normalized before querying to improve consistency.

### Fixed

* Fixed a regression where technicians were incorrectly restricted to viewing only assigned jobs.

### Notes

* MongoDB ObjectIDs continue to be used internally for database relationships.
* Job Numbers are now the preferred external identifier for job retrieval.
