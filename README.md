# README

# Jurassic Park Cage Management Service

This project is designed to allow park personnel to manage dinosaur cage populations and track individual dinosaurs.

## Clone the repository

```
git clone git@github.com:hayesmp/jurassic_cage_service.git
```

## Setting up the Golang app

You will need to have a recent version of the Goland SDK running (this project was written in 1.20). Be sure to download the version appropriate for your processor (arm64 if you're on an Mx Mac or Arm device).
Download SDK: https://go.dev/dl/

In the project directory, run:

```
go mod tidy
```

Assuming no hiccups, you should be ready to start the server or run the tests.


## Running the tests

Simply run all the tests with:

```
cd service/
go test
```

## Running the server

Start the server with:

```
<from the root directory>
go run main.go
```

# Using the API

Once running, there are eight endpoints available to the user:

```
	r.POST("/cage", s.CreateCage)
	r.GET("/cage", s.GetAllCages)
	r.GET("/cage/:id", s.GetCage)
	r.PUT("/cage/:id", s.UpdateCage)
	r.GET("/dinosaur/:id", s.GetDinosaur)
	r.POST("dinosaur", s.CreateDinosaur)
	r.GET("/dinosaur", s.GetAllDinosaurs)
	r.PUT("/dinosaur/:id/:cage_id", s.AddDinosaurToCage)

```

# Cage Endpoints

## Get Cage 

```
GET http://localhost:8080/e00a9b7a-c996-4621-ad15-adb03d8529c5
```

Sample output 200 Success:
```
{
    "id": "e00a9b7a-c996-4621-ad15-adb03d8529c5",
    "name": "Cage 11a761bf-bd81-4195-a761-5574bbe047b0",
    "status": "DOWN",
    "dinosaurs": null,
    "predominate_eating_habit": "Unknown",
    "capacity": 0
}
```

Sample output 404 Not Found

```
{
    "error": "failed to retrieve cage from local db"
}
```


## Create Cage

```
POST http://localhost:8080/cage
{
    "name":"not dangerous dinosaurs"
}
```

Sample 200 Success:
```
{
    "id": "57317bfe-6d28-4257-b4cd-9590df5647aa",
    "name": "not dangerous dinosaurs",
    "status": "ACTIVE",
    "dinosaurs": null,
    "predominate_eating_habit": "Unknown",
    "capacity": 0
}
```


## Get All Cages

```
GET http://localhost:8080/cage
```

Sample 200 Success:
```
[
    {
        "id": "57317bfe-6d28-4257-b4cd-9590df5647aa",
        "name": "not dangerous dinosaurs2",
        "status": "ACTIVE",
        "dinosaurs": null,
        "predominate_eating_habit": "Unknown",
        "capacity": 0
    },
    {
        "id": "69e385a7-3f6d-4742-aa56-0cf8882fa56c",
        "name": "Cage f120a50d-ef33-4940-ad47-2a6bca61b303",
        "status": "ACTIVE",
        "dinosaurs": [
            {
                "id": "ec924278-5b04-4023-bae9-d87e0ab21925",
                "name": "Dino bbafc25e-91c5-41ef-81be-0828bb2ca422",
                "eating_habit": "Herbivore",
                "species": "Stegosaurus",
                "cage_id": "69e385a7-3f6d-4742-aa56-0cf8882fa56c",
                "cage_name": "Cage f120a50d-ef33-4940-ad47-2a6bca61b303"
            },
            {
                "id": "daa15bb8-9102-4987-9f44-5d84e7768201",
                "name": "Dino fc6a68b2-8df3-4a87-a3d8-d9eee35f7cf8",
                "eating_habit": "Herbivore",
                "species": "Ankylosaurus",
                "cage_id": "69e385a7-3f6d-4742-aa56-0cf8882fa56c",
                "cage_name": "Cage f120a50d-ef33-4940-ad47-2a6bca61b303"
            },
            {
                "id": "594ed7cc-f3d2-4c33-979f-fd389c9081fe",
                "name": "Dino 24644c44-7e06-4844-87c1-f3eb17022abc",
                "eating_habit": "Herbivore",
                "species": "Brachiosaurus",
                "cage_id": "69e385a7-3f6d-4742-aa56-0cf8882fa56c",
                "cage_name": "Cage f120a50d-ef33-4940-ad47-2a6bca61b303"
            },
            {
                "id": "6bccd417-7d74-4dd4-b76d-c757225a950d",
                "name": "Dino 74ffea73-f598-49eb-a0bd-99f04b9bd991",
                "eating_habit": "Herbivore",
                "species": "Triceratops",
                "cage_id": "69e385a7-3f6d-4742-aa56-0cf8882fa56c",
                "cage_name": "Cage f120a50d-ef33-4940-ad47-2a6bca61b303"
            }
        ],
        "predominate_eating_habit": "Herbivore",
        "capacity": 4
    },
    ....
]
```

Filter by STATUS

```
GET localhost:8080/cage?status=down
```

Sample 200 Response
```
[
    {
        "id": "57317bfe-6d28-4257-b4cd-9590df5647aa",
        "name": "not dangerous dinosaurs2",
        "status": "DOWN",
        "dinosaurs": null,
        "predominate_eating_habit": "Unknown",
        "capacity": 0
    },
    ...
]
```

## Update Cage

```
PUT http://localhost:8080/cage/57317bfe-6d28-4257-b4cd-9590df5647aa
{
    "status": "down"
}
```

Sample 200 Response
```
{
    "id": "3f859306-34b0-4f2b-8d5c-5c8813b202c4",
    "name": "Cage 11a761bf-bd81-4195-a761-5574bbe047b0",
    "status": "DOWN",
    "dinosaurs": null,
    "predominate_eating_habit": "Herbivore",
    "capacity": 0
}
```

Sample 400 Error
```
{
    "error":"cannot set cage to DOWN with 1 dinosaurs in cage"
}    
```

# Dinosaur Endpoints

## Get Dinosaur

```
GET http://localhost:8080/dinosaur/6298d685-77e3-448c-917e-f80ab3ef989a
```

Sample 200 Success
```
{
    "id": "6298d685-77e3-448c-917e-f80ab3ef989a",
    "name": "Dino f22bf827-5348-4013-9f88-0cb86dc29585",
    "eating_habit": "Herbivore",
    "species": "Brachiosaurus",
    "cage_id": "e7c026f4-bc6c-46e6-92ff-8baeec4ba010",
    "cage_name": "Cage 8602aa8c-dae7-4abd-a4ae-abe2653380e5"
}
```

Sample 404 Not found
```
{
    "error": "error retrieving dinosaur from local db"
}
```

## Create Dinosaur

```
POST http://localhost:8080/dinosaur 
{
    "name":"sam",
    "species":"Tyrannosaurus"
}
```

Sample 200 Success
```
{
    "id": "c274db6b-8214-4319-81b9-096228c31b3a",
    "name": "sams",
    "eating_habit": "Carnivore",
    "species": "Tyrannosaurus",
    "cage_id": "00000000-0000-0000-0000-000000000000",
    "cage_name": ""
}
```

## Get All Dinosaurs

```
GET http://localhost:8080/dinosaur
```

Sample 200 Success
```
[
    {
        "id": "b635c8de-3ad0-4ef1-b9b4-388d07cc4033",
        "name": "ted",
        "eating_habit": "Herbivore",
        "species": "Stegosaurus",
        "cage_id": "00000000-0000-0000-0000-000000000000",
        "cage_name": ""
    },
    {
        "id": "8a1b1666-c373-4b1f-aaf9-b43ebd0aa4bc",
        "name": "sam",
        "eating_habit": "Carnivore",
        "species": "Tyrannosaurus",
        "cage_id": "00000000-0000-0000-0000-000000000000",
        "cage_name": ""
    },
    {
        "id": "ec924278-5b04-4023-bae9-d87e0ab21925",
        "name": "Dino bbafc25e-91c5-41ef-81be-0828bb2ca422",
        "eating_habit": "Herbivore",
        "species": "Stegosaurus",
        "cage_id": "69e385a7-3f6d-4742-aa56-0cf8882fa56c",
        "cage_name": "Cage f120a50d-ef33-4940-ad47-2a6bca61b303"
    },
    ...
]
```

Filter by SPECIES

```
GET http://localhost:8080/dinosaur?species=stegosaurus
```

Sample 200 Success
```
[
    {
        "id": "b635c8de-3ad0-4ef1-b9b4-388d07cc4033",
        "name": "ted",
        "eating_habit": "Herbivore",
        "species": "Stegosaurus",
        "cage_id": "00000000-0000-0000-0000-000000000000",
        "cage_name": ""
    },
    {
        "id": "ec924278-5b04-4023-bae9-d87e0ab21925",
        "name": "Dino bbafc25e-91c5-41ef-81be-0828bb2ca422",
        "eating_habit": "Herbivore",
        "species": "Stegosaurus",
        "cage_id": "69e385a7-3f6d-4742-aa56-0cf8882fa56c",
        "cage_name": "Cage f120a50d-ef33-4940-ad47-2a6bca61b303"
    },
    {
        "id": "6f7027f1-9f3b-49ba-99e5-827831ea704e",
        "name": "Dino f68f22d7-1543-4799-b8f9-497dbae019cd",
        "eating_habit": "Herbivore",
        "species": "Stegosaurus",
        "cage_id": "44d0da26-00a4-4346-8ce7-d28e772643e1",
        "cage_name": "Cage 6c75ddb7-c223-4079-bbbf-49fd509afb3a"
    },
    ...
]
```

## Add Dinosaur To Cage

```
PUT http://localhost:8080/dinosaur/497a8942-d225-43c0-836f-d700a571f6e6/dd3d707a-4649-4771-bd11-7601d788bad3
```

Sample 200 Success
```
{
    "id": "497a8942-d225-43c0-836f-d700a571f6e6",
    "name": "bill",
    "eating_habit": "Herbivore",
    "species": "Stegosaurus",
    "cage_id": "dd3d707a-4649-4771-bd11-7601d788bad3",
    "cage_name": "not dangerous dinosaurs"
}
```

Sample Eating Habit 400 Error
```
{
    "error":"bill cannot be added to cage dd3d707a-4649-4771-bd11-7601d788bad3 due to conflict eating habit"
}
```

Sample Full Cage 400 Error
```
{
    "error":"cage dinosaur count is already at max capacity"
}
```

Sample Cage Status 400 Error
```
{
    "error":"cage Status is not in an ACTIVE state"
}
```

# TODO
- Would be nice to have API Swagger Docs