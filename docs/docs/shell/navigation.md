# Navigation Structure

When applications register themselves with the Gateway they provide information about the entry point, routing, category
and more as defined below. This information is used to construct a navigation structure within the shell and creates the
navigation routes that trigger the loading and rendering of these federated applications.

### Categories

- **App**: places navigation under a root `Apps` menu, this menu supports a depth of 2
      - /myApp
        - /myApp/section1
        - /myApp/section2
- **Settings**: places navigation under a root `Settings` menu, this menu supports a depth of 2
      - /myApp
        - /myApp/section1
        - /myApp/section2
- **Dashboard**: places navigation under a root `Dashboards` menu, this menu supports a depth of 1 and is a mega menu
  - **Slot**: places navigation as Icons within the right side of the site header, max depth is 0 and only pre-defined slots can be used
    - *Note*: There are a maximum of 8 slots (TBD)
    - *Note*: Slots 1 - 5 are reserved for verathread applications and the remaining 3 slots are contextual to the app being loaded
    - Slots are named like `slot-1` and their number denotes the order they will appear in from left to right

### Data Structure

``` json linenums="1"
{
  "defaultRouteId": "34626566-3737-4263-b362-393931383330", // (1)
  "categories": [
    {
      "title": "Apps", // (2)
      "priority": 1, // (3)
      "category": "App",
      "entries": [
        {
          "id": "34626566-3737-4263-b362-393931383330", // (4)
          "title": "Example App Root", // (5)
          "subTitle": "", // (6)
          "icon": "", // (7)
          "authRequired": true, // (8)
          "healthy": true, // (9)
          "module": {
            "moduleName": "ExampleModule", // (10)
            "exposedModule": "./AppModule", // (11)
            "remoteEntry": "/module/34626566-3737-4263-b362-393931383330/remoteEntry.js", // (12)
            "outlet": "", // (13)
            "path": "/example" // (14)
          },
          "children": [ // (15)
            
          ]
        }
      ]
    },
    {
      "title": "Settings",
      "priority": 100,
      "category": "Setting",
      "entries": []
    }
  ],
  "slots": { // (16)
    "slot-1": {
      "description": "Search", 
    },
    "slot-2": {
      "description": "Alerts",
    },
    "slot-3": {
      "description": "Notifications",
    },
    "slot-4": {
      "description": "User", // (17)
      "module": {
        "id": "34626566-3737-4263-b362-393931383330", // (18)
        "authRequired": true,  // (19)
        "healthy": true, // (20)
        "module": {
          "moduleName": "ExampleModule",
          "exposedModule": "./UserModule",
          "remoteEntry": "/module/34626566-3737-4263-b362-393931383331/remoteEntry.js",
          "path": "/user"
        }
      }
    }
  }
}

```

1.  :fontawesome-solid-circle-info: The id of the default route to navigate to when the path is `/`
2.  :fontawesome-solid-circle-info: The title of the root navigation item in the header
3.  :fontawesome-solid-circle-info: The order this menu will appear in the header, lower priority will appear first from left to right
4.  :fontawesome-solid-circle-info: The unique id of the navigation entry
5.  :fontawesome-solid-circle-info: The menu title, this is a short value, e.g. Users, Account, Rune
6.  :fontawesome-solid-circle-info: The menu subtitle, displayed inline if possible otherwise a tooltip, this should be a short and direct description of the app
7.  :fontawesome-solid-circle-info: A base64 encoded svg icon or font-awesome icon name e.g. `fa-info`
8.  :fontawesome-solid-circle-info: If true then the route for the nav entry should have an auth guard attached
9.  :fontawesome-solid-circle-info: If true the nav item should be enabled, otherwise should be disabled with a tooltip explaining that the service is currently unavailable
10.  :fontawesome-solid-circle-info: The remote entry module name to use during routing
11.  :fontawesome-solid-circle-info: The exposed module name as defined in the webpack federation configuration
12.  :fontawesome-solid-circle-info: The url to the remote entry file, used by the router to load the entry point of a federated app
13.  :fontawesome-solid-circle-info: Outlet, only used if the target of the navigation uses named outlets, the router can load the module into that outlet 
14.  :fontawesome-solid-circle-info: Path, the base path for deep linking, the rest of the path is populated by the app's router (must be globally unique otherwise your application will fail to register)
15.  :fontawesome-solid-circle-info: Sub paths, these can not load modules but instead allow navigating to a sub path of the app directly from the root menu `/root/child` (To be supported later)
16.  :fontawesome-solid-circle-info: Slots are named by the order they would appear on the right side of the header, apps making use of slots must make sure their app renders an icon and not content suited to a page
17.  :fontawesome-solid-circle-info: Description of the slot item, can be used for logging or error handling
18.  :fontawesome-solid-circle-info: The id of the module
19.  :fontawesome-solid-circle-info: If true then the icon should be disabled with a tooltip letting the user know that they need to log in first
20.  :fontawesome-solid-circle-info: If true then the slot should be enabled, otherwise disabled with a tooltip letting the user know that the service is unavailable



