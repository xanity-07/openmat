# OpenMat Authentication Roadmap

This document defines the authentication and account-security features for OpenMat. The goal is to build a secure, production-minded authentication system without introducing unnecessary complexity.

---

## 1. Registration

### 1.1 Define the User Model

- [ ] Create `users` table
- [ ] Generate UUID user IDs
- [ ] Add `email`
- [ ] Add `password_hash`
- [ ] Add `created_at`
- [ ] Add `updated_at`
- [ ] Add `deleted_at` if using soft deletion
- [ ] Add appropriate indexes and constraints

### 1.2 Input Validation

- [ ] Validate email format
- [ ] Normalize email where appropriate
- [ ] Validate password requirements
- [ ] Validate request payload
- [ ] Return consistent validation errors

### 1.3 Password Hashing

- [ ] Choose password hashing algorithm
- [ ] Hash passwords before persistence
- [ ] Never store plaintext passwords
- [ ] Never log passwords
- [ ] Verify hashes using the hashing library's verification function

### 1.4 Email Uniqueness

- [ ] Add database-level unique constraint
- [ ] Handle duplicate-email database errors
- [ ] Return a consistent API error

### 1.5 Registration Flow

```text
Request
  ↓
Validate payload
  ↓
Normalize email
  ↓
Check/create user
  ↓
Hash password
  ↓
Persist user
  ↓
Create session
  ↓
Issue tokens
  ↓
Return authentication response
```

---

# 2. Login

### 2.1 Credential Verification

- [ ] Find user by email
- [ ] Verify password hash
- [ ] Reject deleted/inactive accounts
- [ ] Return a generic authentication failure
- [ ] Avoid revealing whether an email exists

### 2.2 Access Token

- [ ] Define JWT claims
- [ ] Include user/session identifier
- [ ] Set expiration
- [ ] Sign with a strong secret/key
- [ ] Validate issuer/audience if applicable
- [ ] Validate signing algorithm when parsing tokens

### 2.3 Refresh / Session Token

- [ ] Decide whether OpenMat needs refresh tokens
- [ ] Generate cryptographically secure random session tokens
- [ ] Store only a hash of the session token if practical
- [ ] Associate sessions with users
- [ ] Set session expiration
- [ ] Support session revocation

### 2.4 Login Flow

```text
Email + Password
       ↓
Validate input
       ↓
Find user
       ↓
Verify password
       ↓
Create session
       ↓
Issue access token
       ↓
Issue refresh/session token
```

---

# 3. Sessions

Redis can be used for fast session/revocation state while PostgreSQL remains the source of persistent user/account data.

### 3.1 Redis-backed Sessions

- [ ] Define session structure
- [ ] Generate unique session IDs
- [ ] Associate session with user ID
- [ ] Store session expiration
- [ ] Store relevant metadata
- [ ] Use Redis TTLs

Example conceptual structure:

```text
session:{session_id}
    user_id
    created_at
    expires_at
```

### 3.2 Session Expiration

- [ ] Define session lifetime
- [ ] Set Redis TTL
- [ ] Reject expired sessions
- [ ] Remove expired session state

### 3.3 Logout / Revocation

- [ ] Authenticate current session
- [ ] Delete/revoke session
- [ ] Invalidate refresh token
- [ ] Prevent revoked sessions from being refreshed

### 3.4 Logout All Sessions

- [ ] Find all sessions belonging to the user
- [ ] Revoke all active sessions
- [ ] Ensure current and other devices are invalidated
- [ ] Consider a user-level session/version mechanism to make this efficient

### 3.5 Session Management

Eventually consider:

- [ ] List active sessions
- [ ] Revoke individual session
- [ ] Track device/user-agent information
- [ ] Track IP information where appropriate
- [ ] Display session creation/last-used timestamps

---

# 4. Authorization

Authentication answers:

> Who are you?

Authorization answers:

> Are you allowed to do this?

### 4.1 Authentication Middleware

- [ ] Extract access token
- [ ] Validate token
- [ ] Validate expiration
- [ ] Validate signing method
- [ ] Resolve user/session
- [ ] Reject unauthenticated requests
- [ ] Attach authenticated identity to request context

### 4.2 User Identity / Context

- [ ] Define authenticated-user context
- [ ] Make user ID available to handlers
- [ ] Avoid trusting user IDs supplied by clients
- [ ] Derive ownership from authenticated identity

Example:

```text
JWT
 ↓
Authentication middleware
 ↓
User ID
 ↓
Request context
 ↓
Handler
 ↓
Repository query scoped to authenticated user
```

### 4.3 Roles / Permissions

Only implement this if OpenMat actually needs it.

- [ ] Define roles
- [ ] Define permissions
- [ ] Add authorization middleware/helpers
- [ ] Test unauthorized access
- [ ] Test cross-user resource access

---

# 5. Account Security

## 5.1 Email Verification

- [ ] Create email-verification token
- [ ] Generate cryptographically secure token
- [ ] Store token securely
- [ ] Set expiration
- [ ] Send verification email
- [ ] Verify token
- [ ] Mark email as verified
- [ ] Invalidate token after use
- [ ] Support requesting another verification email

## 5.2 Password Reset

- [ ] Create password-reset token
- [ ] Generate cryptographically secure token
- [ ] Store token securely
- [ ] Set short expiration
- [ ] Send reset email
- [ ] Validate token
- [ ] Allow password replacement
- [ ] Invalidate token after use
- [ ] Invalidate existing sessions after successful reset

**Important:** Password-reset tokens should not be reusable.

## 5.3 Password Change

- [ ] Require authenticated session
- [ ] Require current password
- [ ] Validate new password
- [ ] Hash new password
- [ ] Replace password hash
- [ ] Invalidate appropriate sessions
- [ ] Require reauthentication where appropriate

## 5.4 Session Invalidation After Password Changes

When credentials change:

```text
Password changed
      ↓
Invalidate existing sessions
      ↓
Existing refresh tokens stop working
      ↓
User must authenticate again
```

This prevents an attacker who obtained a user's session from retaining access after the legitimate user changes their password.

---

# 6. Security Controls

## 6.1 Rate Limiting

Prioritize authentication endpoints.

- [ ] Rate-limit login attempts
- [ ] Rate-limit registration
- [ ] Rate-limit password-reset requests
- [ ] Rate-limit verification-email requests
- [ ] Decide whether limits should be IP-based, account-based, or both
- [ ] Use Redis if appropriate
- [ ] Return consistent responses

Be careful not to create an account-enumeration vulnerability through different rate-limit behavior.

## 6.2 Secure Cookies

If authentication tokens are stored in cookies:

- [ ] `HttpOnly`
- [ ] `Secure` in production
- [ ] Appropriate `SameSite`
- [ ] Appropriate cookie expiration
- [ ] Restrict cookie scope where possible

Do not store sensitive authentication tokens in ordinary client-accessible storage without a deliberate security reason.

## 6.3 CSRF Protection

Determine whether CSRF protection is required based on the authentication mechanism.

- [ ] Document the decision
- [ ] Implement CSRF protection if using cookie-based authentication
- [ ] Validate CSRF tokens where required
- [ ] Test cross-site requests

If authentication uses cookies, treat CSRF as a real concern rather than assuming CORS solves it.

## 6.4 Consistent Authentication Errors

Define a consistent error model.

Examples:

```text
401 Unauthorized
403 Forbidden
422 Validation Error
429 Too Many Requests
```

Avoid responses that reveal sensitive information.

For example, prefer:

```text
Invalid email or password
```

over:

```text
No account exists with this email
```

## 6.5 Audit / Security Logging

Log security-relevant events without logging secrets.

Potential events:

- [ ] Login succeeded
- [ ] Login failed
- [ ] Logout
- [ ] Session revoked
- [ ] All sessions revoked
- [ ] Password changed
- [ ] Password reset requested
- [ ] Password reset completed
- [ ] Email verified
- [ ] Account created
- [ ] Suspicious authentication activity

Never log:

- passwords
- JWTs
- refresh tokens
- session tokens
- password-reset tokens
- verification tokens
- other authentication secrets

---

# 7. Testing

Authentication should have significantly more test coverage than ordinary CRUD code.

### Registration

- [ ] Valid registration
- [ ] Invalid email
- [ ] Weak password
- [ ] Duplicate email
- [ ] Deleted user/email behavior

### Login

- [ ] Valid credentials
- [ ] Invalid password
- [ ] Unknown email
- [ ] Deleted account
- [ ] Expired session
- [ ] Invalid token
- [ ] Invalid signing algorithm

### Sessions

- [ ] Session creation
- [ ] Session expiration
- [ ] Logout
- [ ] Session revocation
- [ ] Logout all sessions
- [ ] Refresh with revoked session

### Authorization

- [ ] Unauthenticated request
- [ ] Authenticated request
- [ ] User accessing own resource
- [ ] User attempting to access another user's resource
- [ ] Unauthorized role/permission

### Account Security

- [ ] Email verification
- [ ] Expired verification token
- [ ] Password reset
- [ ] Expired reset token
- [ ] Reused reset token
- [ ] Password change
- [ ] Session invalidation after password change

### Rate Limiting

- [ ] Requests below limit succeed
- [ ] Requests above limit are rejected
- [ ] Limits expire/reset correctly

---

# 8. Implementation Order

Don't build everything simultaneously.

### Phase 1 — Core Authentication

- [ ] User database model
- [ ] Registration
- [ ] Password hashing
- [ ] Login
- [ ] Access tokens
- [ ] Authentication middleware

### Phase 2 — Sessions

- [ ] Redis session storage
- [ ] Session expiration
- [ ] Logout
- [ ] Session revocation
- [ ] Refresh/session tokens
- [ ] Logout all sessions

### Phase 3 — Account Recovery

- [ ] Email verification
- [ ] Password reset
- [ ] Password change
- [ ] Session invalidation after credential changes

### Phase 4 — Security Hardening

- [ ] Rate limiting
- [ ] Secure cookies
- [ ] CSRF protection where applicable
- [ ] Consistent authentication errors
- [ ] Security logging

### Phase 5 — Testing & Review

- [ ] Unit tests
- [ ] Integration tests
- [ ] Authentication middleware tests
- [ ] Session life-cycle tests
- [ ] Security-focused tests
- [ ] Review token/session expiration
- [ ] Review secrets/configuration
- [ ] Review logs for accidental credential exposure

---

# Definition of Done

OpenMat authentication should be considered production-minded when:

- Passwords are securely hashed.
- Authentication tokens have controlled lifetimes.
- Sessions can be revoked.
- Logout actually invalidates the session.
- Users can invalidate all sessions.
- Password changes/reset invalidate existing sessions appropriately.
- Authentication endpoints are rate-limited.
- Sensitive authentication material is never logged.
- Authentication and authorization are tested.
- Users cannot access another user's resources by manipulating IDs.
- Authentication failures do not unnecessarily disclose account information.
- Secrets are provided through configuration rather than committed to source control.
