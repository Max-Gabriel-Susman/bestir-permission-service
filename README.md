# Bestir permissionMaking Service 

Handles the lifecycle of permissions on the bestir network

### Roadmap 

Let's start w/ rudimentary permission creation and management, testing for that logic, I'll build the service to be production grade and then scale up operations from there

to get the bestir platform working with the unreal game templates I think we're going to initially just need an identity, permissions, and applicaton(game creation and management) service

let's work out the simplest version of this system to allow the hosting of the linux images for the unreal template games, and then build from there

## dev notes

[11/27/2022] included delete and update logic for permission endpoints, but rn it's basically just copypasta of create logic so def not ready for use there