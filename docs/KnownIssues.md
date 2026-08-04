# Known Issues

## React Router Advisory

**Package**
- react-router-dom

**Advisory**
- GHSA-qwww-vcr4-c8h2

**Severity**
- High

**Status**
- Accepted Risk

**Reason**

This project is a React SPA using:

- Vite
- Gin REST API
- JWT Authentication
- Client-side routing

The reported vulnerability affects **React Server Components (RSC)** and server actions.

This application does **not** use:

- React Server Components
- React Router server actions
- React Router loaders on the server

Therefore the advisory is currently **not applicable**.

The dependency should be reviewed again when React Router publishes an updated release or the advisory database is corrected.
