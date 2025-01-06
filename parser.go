/*
 * Author: Nathaniel Martes
 * Description: Parses PGN file and stores data in a database described in chess.sql
 */


package main

import (
	"os"
	"github.com/notnil/chess"
    "fmt"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
	"strings"
	"time"
)

/*
 * MoveData is used for preparing moves to be inserted into the Moves MySQL table
 */
type MoveData struct {
	StartPos string
	EndingPos string
	Color string
	MoveNum int
	KingSideCastle int
	QueenSideCastle int
}
func main() {

	//open pgn file
    chessFile, err := os.Open("twic210-874.pgn")
	if err != nil {
		panic(err)
	}
	defer chessFile.Close()

	scanner := chess.NewScanner(chessFile)

	//Open Database and check connection
	db, err := sql.Open("mysql", "root:CS235@tcp(mysql:3306)/chess")
	if err != nil {
		panic(err.Error())
	}
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	insertEventQuery := "INSERT INTO events (event_id, name, site, date, round) VALUES (?, ?, ?, ?, ?)"
	insertGamesQuery := "INSERT INTO games (game_id, event_id, white, black, result) VALUES (?, ?, ?, ?, ?)"
	insertMovesQuery := "INSERT INTO moves (game_id, move_num, color, start_pos, ending_pos, kingside_castl, queenside_castl) VALUES %s"

	table_id := 1;

	var errArr []string //needed to store errors that occur during insertion
	failedEvent := "Failed to store event at old event_id : %d [%s]"
	failedGame := "Failed to store game at event_id : %d [%s]"
	failedMoves := "Failed to store moves at game_id : %d [%s]"

	for scanner.Scan() {
		game := scanner.Next()
		
		event := game.GetTagPair("Event").Value
		site := game.GetTagPair("Site").Value
		date := game.GetTagPair("Date").Value
		white := game.GetTagPair("White").Value
		black := game.GetTagPair("Black").Value
		round := game.GetTagPair("Round").Value
		result := game.GetTagPair("Black").Value

		//Change result to 1 character
		if result == "1/2-1/2" {
			result = "T"
		} else if result == "1-0" {
			result = "W"
		} else {
			result = "B"
		}

		//Parse move data
		moves := make([]MoveData, len(game.Moves()))
		gameMoves := game.Moves()

		i := 0
		j := 1
		moveNum := 1

		//since even indexes are white moves and odd indexes are black moves, we can add 2 moves per iteration
		for i < len(gameMoves) {
			KSC := 0
			QSC := 0
			if gameMoves[i].HasTag(chess.KingSideCastle) {
				KSC = 1
			}
			if gameMoves[i].HasTag(chess.QueenSideCastle) {
				QSC = 1
			}
			moves[i] = MoveData{
				StartPos: gameMoves[i].S1().String(),
				EndingPos: gameMoves[i].S2().String(),
				Color: "W",
				MoveNum: moveNum,
				KingSideCastle: KSC,
				QueenSideCastle: QSC,
			}

			//incase the final move only has 1 white move
			if j >= len(gameMoves) {
				break
			}

			KSC = 0
			QSC = 0
			if gameMoves[j].HasTag(chess.KingSideCastle) {
				KSC = 1
			}
			if gameMoves[j].HasTag(chess.QueenSideCastle) {
				QSC = 1
			}
			moves[j] = MoveData{
				StartPos: gameMoves[j].S1().String(),
				EndingPos: gameMoves[j].S2().String(),
				Color: "B",
				MoveNum: moveNum,
				KingSideCastle: KSC,
				QueenSideCastle: QSC,
			}
			i += 2
			j += 2
			moveNum++
		}

		/*
		 * Inserting Data into Tables
		 * Note: 3 cases exists when inserting data into tables
		 *	1. Event data is bad: All data for that chess game is skipped
		 *	2. Game data is bad: Event data is saved but game and move data is skipped
		 *  3. Move data is bad: Event and game data is saved but move data is skipped
		 */


		//Check if round number and/or date exist
		if round == "?" {
			if strings.Contains(date,"?") {
				_, err = db.Exec(insertEventQuery, table_id, event, site, nil, nil)
			} else {
				_, err = db.Exec(insertEventQuery, table_id, event, site, date, nil)
			}
		} else {
			if strings.Contains(date,"?") {
				_, err = db.Exec(insertEventQuery, table_id, event, site, nil, round)
			} else {
				_, err = db.Exec(insertEventQuery, table_id, event, site, date, round)
			}
		}
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Skipping event, game, and moves for event : "+event+" "+site)
			errArr = append(errArr, fmt.Sprintf(failedEvent, table_id, err.Error()))
			continue
		} else {
			fmt.Println("Inserted new event " + event)
		}
		

		_, err = db.Exec(insertGamesQuery, table_id, table_id, white, black, result)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Printf("Skipping game, and moves for event_id : %d\n", table_id)
			errArr = append(errArr, fmt.Sprintf(failedGame, table_id, err.Error()))
			table_id++
			continue
		} else {
			fmt.Println("Inserted new game " + white + " VS. " + black)
		}

		//since moves for a game may be empty somehow, we will check just in case
		if len(moves) != 0 {
			//setup moves for batch insert to reduce time
			var valueStrings []string
			var valueArgs []interface{}
			for _, move := range moves {
				valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?, ?)")
				valueArgs = append(valueArgs, fmt.Sprintf("%d", table_id,), fmt.Sprintf("%d", move.MoveNum),fmt.Sprintf("%s", move.Color), fmt.Sprintf("%s", move.StartPos),
				fmt.Sprintf("%s", move.EndingPos), fmt.Sprintf("%d", move.KingSideCastle), fmt.Sprintf("%d", move.QueenSideCastle))
			}

			//perform batch insert on moves
			stmt := fmt.Sprintf(insertMovesQuery, strings.Join(valueStrings, ","))
			_, err = db.Exec(stmt, valueArgs...)
			if err != nil {
				fmt.Println(err.Error())
				fmt.Printf("Skipping moves for game_id : %d\n", table_id)
				errArr = append(errArr, fmt.Sprintf(failedMoves, table_id, err.Error()))
				table_id++
				continue
			}
			fmt.Println("Inserted moves for " + white + " VS. " + black)
		}
		table_id++
	}
	//Print all errors that occured:
	time.Sleep(2 * time.Second)
	fmt.Println("========================")
	fmt.Println()
	fmt.Printf("Insertion Errror Count : %d\n", len(errArr))
	fmt.Println()
	for _, err := range errArr {
		fmt.Println(err)
	}
	fmt.Println()

	time.Sleep(2 * time.Second)


	//Execute Test Queries
	fmt.Println("Executing Test Queries")
	fmt.Println()
	time.Sleep(1 * time.Second)
	fmt.Println("How many total games did black win?")
	fmt.Println("Executing Query : SELECT COUNT(*) FROM games WHERE result = \"B\"")
	time.Sleep(1 * time.Second)

	rows, err := db.Query("SELECT COUNT(*) FROM games WHERE result = ?", "B")
    if err != nil {
        panic(err)
    }
	for rows.Next() {
        var result int
        err := rows.Scan(&result)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Result : %d\n", result)
    }
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)




	fmt.Println("What percentage of games in the database open with b4?")
	fmt.Println("Executing Query : SELECT (COUNT(*) / (SELECT COUNT(*) FROM moves WHERE move_num = 1 AND color = \"W\"))*100 as \"% of b4 opening games\" FROM moves WHERE move_num = 1 and ending_pos = \"b4\";")
	time.Sleep(1 * time.Second)

	rows, err = db.Query("SELECT (COUNT(*) / (SELECT COUNT(*) FROM moves WHERE move_num = 1 AND color = ?))*100 FROM moves WHERE move_num = 1 and ending_pos = ?", "W", "b4")
    if err != nil {
        panic(err)
    }
	for rows.Next() {
        var result float32
        err := rows.Scan(&result)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Result : %f\n", result)
    }
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	
	fmt.Println("What is the average length of a chess game (in total moves)?")
	fmt.Println("Executing Query : SELECT ROUND(SUM(max_move_in_game)/(SELECT COUNT(*) FROM games WHERE game_id IN (SELECT game_id FROM moves)),0) as \"avg moves per game\" FROM (SELECT MAX(move_num) as max_move_in_game FROM moves GROUP BY game_id) as q;")
	time.Sleep(1 * time.Second)

	rows, err = db.Query("SELECT ROUND(SUM(max_move_in_game)/(SELECT COUNT(*) FROM games WHERE game_id IN (SELECT game_id FROM moves)),0) FROM (SELECT MAX(move_num) as max_move_in_game FROM moves GROUP BY game_id) as q;")
    if err != nil {
        panic(err)
    }
	for rows.Next() {
        var result float32
        err := rows.Scan(&result)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Result : %f\n", result)
    }
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)

	fmt.Println("How many total games include kingside castling on or before black's twentieth move?")
	fmt.Println("Executing Query : SELECT COUNT(DISTINCT game_id) FROM moves WHERE kingside_castl = 1 AND move_num <= 20;")
	time.Sleep(1 * time.Second)

	rows, err = db.Query("SELECT COUNT(DISTINCT game_id) FROM moves WHERE kingside_castl = 1 AND move_num <= 20;")
    if err != nil {
        panic(err)
    }
	for rows.Next() {
        var result int
        err := rows.Scan(&result)
        if err != nil {
            panic(err)
        }
        fmt.Printf("Result : %d\n", result)
    }
	fmt.Println()
	fmt.Println()
	time.Sleep(2 * time.Second)
}