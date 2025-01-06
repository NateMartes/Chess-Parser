# Chess-Parser
A PGN file parser that inserts data into a MySQL Database.
This was a project for my CS235 class so you'll notice some tid bits throughout the files :)

## Running the Parser
This project is only meant to be ran through Docker and Docker Compose so they will be needed:
 - Docker : https://docs.docker.com/engine/install/
 - Docker Compose : https://docs.docker.com/compose/install/

Things to Note:

 - This project takes about 9 minutes to store the data from the pgn file even with using batch inserts.
 - If you have Docker containers using ports 8080 and 3306, there may be an issue build containers so just pause any if they exist OR change the ports being used in docker-compose.yml AND parser.go .

1. Clone the repository.

	`git clone https://github.com/NateMartes/Chess-Parser.git`

2. Add your PGN file to the repository (Example: MyPGNFile.pgn).
3. Change line 33 in parser.go.

	```chessFile, err := os.Open("INSERT PGN FILE HERE")```

4. Run Docker Compose in the current repository.

	```sudo docker-compose up --build```

## Resources Used

[Chess Library](https://github.com/notnil/chess): A set of go packages which provide common chess utilities such as move generation, turn management, checkmate detection, PGN encoding, UCI interoperability, image generation, opening book exploration, and others. It is well tested and optimized for performance.

[wait-for-it](https://github.com/vishnubob/wait-for-it): A pure bash script that will wait on the availability of a host and TCP port. It is useful for synchronizing the spin-up of interdependent services, such as linked docker containers. Since it is a pure bash script, it does not have any external dependencies.


