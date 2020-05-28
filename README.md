# ultimate-jwt-auth-server
A JWT authentication server implementation using JWT access tokens and refresh tokens

Inspired by this Hasura article: https://hasura.io/blog/best-practices-of-using-jwt-with-graphql/#intro
Uses JWT to authenticate users.

This Golang authentication server supports a partially decentralized approach to session management via JWT authentication.

## Workflow

Upon registering via /register, users are given two items:
- a short-lived JWT access token with expiry, which can be sent to other endpoints (like an API) to authenticate via JWS (no need to keep a centralized session DB).
- a long-lived token sent as an HttpOnly cookie with scope set to the /refresh endpoint of the auth server

Users can access other cooperating endpoints (which correctly implement JWT logic) without having to sign in. These cooperating servers also do not need to save session state since JWT is decentralized by design.

Upon expiration of the short-lived JWT access tokens, users will be redirected to the /refresh endpoint of the auth server. 
Here, the auth server will verify that the refresh token cookie is valid and not expired before regenerating new access and refresh tokens.
If the refresh token is expired or invalid, users will be redirected to /login.

At /login, if user sends correct username/email + password data, auth server will generate new access and refresh tokens.

## Database

The server is hooked up to a local instance of PostgreSQL, which is interfaced with the popular Golang Gorm ORM.
Database and server can be configured out via .env variables.
