# Authentication Flow Design

## Objective

The primary goal for authentication was to implement a secure, standard-based OAuth 2.0 flow suitable for a desktop CLI application.

## Selected Flow: Authorization Code Flow for Desktop Apps

After extensive testing and research, the application uses Google's standard **Authorization Code Flow** for installed applications.

### The Role of the Client Secret

A key aspect of this implementation is the use of a `client_secret` that is embedded in the compiled binary. This approach can seem counter-intuitive based on general OAuth 2.0 best practices, which advise against storing secrets in public clients.

However, Google's specific implementation for "Desktop app" credentials differs from the general standard. For this client type, Google provides a `client_secret` and requires it to be sent during the token exchange. This secret is not treated as a highly confidential value in the same way a web server's secret is. Instead, it acts as a stable identifier for the application.

This design choice was made based on the following:
1.  **Persistent Errors with PKCE**: Attempts to use the PKCE (Proof Key for Code Exchange) flow without a client secret consistently failed with a `client_secret is missing` error from Google's servers. This occurred even when using a credential explicitly created as a "Desktop app" type.
2.  **Community Confirmation**: Research indicates that this is expected behavior for Google's platform. The `client_secret` for an installed application is required and is intended to be distributed with the application.

### Security Acknowledgment

While embedding a value named `client_secret` in a client application is non-standard for other platforms, it is the required and documented method for Google's OAuth 2.0 flow for desktop applications. The security of the flow is maintained by the use of user-specific authorization codes and refresh tokens, which are stored securely on the user's machine.

### Sources

The decision to proceed with this design was informed by community discussions and documentation that clarify Google's specific requirements for desktop applications.

- [Google's OAuth 2.0 for installed apps and Client Secret not being a secret](https://stackoverflow.com/questions/44312000/googles-oauth-2-0-for-installed-apps-and-client-secret-not-being-a-secret)
- [Safely distribute OAuth 2.0 client_secret in desktop applications in Python](https://stackoverflow.com/questions/59416326/safely-distribute-oauth-2-0-client-secret-in-desktop-applications-in-python)
- [Google Cloud OAuth 2.0 missing client secret](https://stackoverflow.com/questions/78878049/google-cloud-oauth-2-0-missing-client-secret)
- ["error_description": "client_secret is missing." - Google Groups](https://groups.google.com/g/gce-discussion/c/zds-LP15rc8)
- [Is a Client Secret Required for Google OAuth 2.0 using the PKCE authorization flow?](https://stackoverflow.com/questions/78673050/is-a-client-secret-required-for-google-oauth-2-0-using-the-pkce-authorization-fl)