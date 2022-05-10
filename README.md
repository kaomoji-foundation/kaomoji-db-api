# Go API template
## Using/implementing
<details>
 <summary>Click to expand!</summary>
 - **Fiber**    
    > For performant routing and many other things
 - **JWT**      
    > For user authentication
 - **MongoDB**  
    > As the main database
 - **OpenTelemetry & LightStep**  
    > For Deployment telemetry
 - **Redoc standalone** 
    > For Technical documentation deployment
 - **godotenv** 
    > For local development
</details>

## Introduction/vision:

Fiber is a quickly growing super-fast backend framework, mainly written in golang with performance in mind. 

The intent of this repo is to create a **performant, comprenhensible, mantainable and production ready template** with nothing in it but the basics, so you do not have to spend time removing code or strugeling to understand any complicated code structure

The only thing it handles by default is a sesion management system.

## Future features:
<details>
 <summary>Click to expand!</summary>
   ### Nothing to see here yet  :(
</details>

## Project Structure:
<details>
 <summary>Click to expand!</summary>

```
.
├── src/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   └── users/
│   │       ├── auth.go
│   │       └── user.go
│   ├── helpers/
│   │   └── .
│   ├── loaders/
│   │   ├── fiber.go
│   │   └── mongo.go
│   ├── middlewares/
│   │   └── auth.go
│   ├── models/
│   │   └── user-model.go
│   ├── routes/
│   │   ├── appRouter.go
│   │   ├── auth.go
│   │   └── users.go
│   ├── services/
│   │   └── mongo.go
│   ├── utils/
│   │   └── .
│   └── main.go
├── .env
├── app.go
├── go.mod
├── go.sum
├── README.md
└── sample.env
```
</details>

I wanted to contain all the source code under a sub folder, so I only use app.go to execute a Start() function in src/main.go

# API documentation
   currently swaggest/swag cli **was** being used to generate the docs from the decalrative comments, this is no longer the case.
   **Currently using manualy crafted documents**
   > Why not use swaggest/swag cli?: it does not support OAS3

   To serve theese docs a **Redoc** static file is being used, note that internet connection is required as jdelivr is being used for dependencies.

## Setup:
<details>
  <summary>Click to expand!</summary>

To get the API up and running, there are many things to setup first,  

### Building from source:
<details>
  <summary>Click to expand!</summary>
1. Once this repository is cloned and golang is installed in the system, navigate to this directory and run

```
go mod download
```

2. Once the dependencies are downloaded, using sample.env as reference either create a file called creds.env with the same keys or directly configure same keys as environment variables.
3. After the configuration and ensuring that the db is operational, run either

```
go build
```
 to get the executable to run

or to build and run with a single command

```
go run app.go
```

</details>

### MongoDB setup (local)

<details>
  <summary>Click to expand!</summary>
You will need a mongoDB database, with atleast 
   
- 1 user with credentials. 
  - defaults:  (API will use theese values by default)
    - username: `root`
    - password: `example`
- 1 database.
- at least a single role in the roles document.

#### Default recomended roles[ready to import]:
```json
[{
  "_id": {
    "$oid": "626ad6e35204187f3579d44f"
  },
  "role": "admin",
  "level": 1,
  "permissons": {
    "readUsers": true,
    "usersAdmin": true,
    "readRoles": true,
    "rolesAdmin": false
  }
},{
  "_id": {
    "$oid": "626d011e9c6806ef1f5cddd5"
  },
  "role": "user",
  "level": 3,
  "permissons": {
    "readUsers": true,
    "usersAdmin": false,
    "readRoles": false,
    "rolesAdmin": false
  }
},{
  "_id": {
    "$oid": "626d01589c6806ef1f5cddd6"
  },
  "role": "root",
  "level": 0,
  "permissons": {
    "readUsers": true,
    "usersAdmin": true,
    "readRoles": true,
    "rolesAdmin": true
  }
},{
  "_id": {
    "$oid": "626d01cb9c6806ef1f5cddd7"
  },
  "role": "moderator",
  "level": 2,
  "permissons": {
    "readUsers": true,
    "usersAdmin": true,
    "readRoles": false,
    "rolesAdmin": false
  }
}]
``` 

</details>

### Configuring local API

<details>
  <summary>Click to expand!</summary>
   Copy `sample.env` into `.env`
   and fill the `.env`

   For local development if you are using the default
   you should be good to go by just filling:
   - `DB_1_NAME`
   - `JWT_SECRET`

   `DB_1_NAME` being the database name, and `JWT_SECRET` a random string. 
   This last one theoridically is not necessary, but is strongly recomended.

</details>

### Configure OpenTelemetry (Optional)

<details>
  <summary>Click to expand!</summary>
   Open telemetry requires no etra configuration, but the data collector used does, in our case it is [LightStep](https://lightstep.com)

   If you dont like this, dont fear, swaping this is actually very simple, since i have not integrated it 100%, and is only a few lines in `./src/main.go` what you would have to remove/swap:
   
   ```go
   //Open Telemetry setup
	ls := launcher.ConfigureOpentelemetry(
		launcher.WithServiceName("Go-API-Template"),
		launcher.WithAccessToken(cfg.Config.OpenTel.LightStepKey),
	)
	defer ls.Shutdown()
	// END Open Telemetry setup
   ```
   ###### *This may not correspond 100% to reality 
</details>

</details>


## About:

If this repository is/was useful to you in any way, please star this repository and share it with people who may be interested. I'll do my best to keep it updated.
