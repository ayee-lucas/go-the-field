
## Features V2 - Changelog

### Database Schema Enhancement

-   Upgraded the MongoDB database schema to improve reliability and structure.
-   Renamed the "Org" collection and route handlers to "Teams" for better clarity and consistency.

### Improved User Struct

-   Refactored the User struct into two separate entities: Profile and Type.
-   The Profile struct now contains all the necessary data for user profiles.
-   Introduced a 'Type' field in the profile to store the user's type, including team/athlete ids.

### Streamlined Relationship Management

-   Replaced the array-based approach for Likes/Followers/Following properties with a single object id.
-   Likes/Followers/Following documents now utilize this object id as a 'host' to maintain user-document relationships.
-   This enhancement addresses issues with the previous 'unlimited-array' approach, ensuring greater reliability for the collection.

### Bug Fixes and Enhancements

-   Addressed typos and improved error messages for better user experience.
-   Fixed issues related to data types when saving documents to the database.

### Refocused Project Scope

-   Initially designed as the sole backend for "THE FIELD" project, this Golang project's direction has been adjusted.
-   The primary purpose now revolves around authentication management, including SignUp, Login, Sessions, Logout, GetUser, and finishing user profiles.
-   Other functionalities and features will be handled by the NextJs API using Prisma.

##
##
##

## **FeaturesV1**

### Improved Error Messages

- Enhanced the responses returned from the `accountHandler` route handler to provide more informative and user-friendly error messages.

### Signup Function Enhancement

- Added the missing field to the `User` struct in the `accountHandler` signup function, ensuring all necessary user data is captured during signup.

### Route Handler Functionality

- Implemented the `RequestOrg` function to handle specific routes related to organizations.

### User Struct Extension

- Expanded the `User` struct to include additional fields for `Conversations`, `Sports`, `Org`, and `Athlete`, enabling better user data representation.

### User Repository Improvements

- Added several utility functions to the user repository:
  - `Save`: For saving user data.
  - `GetById`: For retrieving user data by user ID.
  - `GetByEmail`: For fetching user data based on the email address.
  - `DeleteOrgById`: For deleting an organization by its ID.

### New Route

- Added the `request/org/:id` route to handle specific requests related to organizations.
