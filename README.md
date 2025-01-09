 # chess-arrangements

White - 216481 options

Black - 216481 options

Together - 46,864,023,361 options
 

To run the program:

```
# For just white and black separate files: 
go run main.go 

# For full combinations: 
go run main.go full
```

The files in the `full` directory will contain:

```
positions_1.txt:
Position 1: RNBQKBNR/PPPPPPPP-rnbqkbnr/pppppppp
Position 2: RNBQKBNR/PPPPPPPP-rnbqkbnr/pppppppq
...
(1 million positions)

positions_2.txt:
Position 1000001: ...
...
```


And a `summary.txt`:

```
Total combinations: 46864037361
Positions per file: 1000000
Total files: 47
```

---

## Prompts

On the first two horizontals the player can place his pieces in any order, namely: 8 pawns, 2 rooks, 2 knights, 2 bishops, 1 queen and 1 king, the record of each possible arrangement is written in a simplified `fen` format, it is necessary to write a program in GO lang that writes all possible options to a file indicating the option number. Pawns and pieces can be placed on both the 1st and 2nd ranks

Do this separately, for both white and black pieces with separate files

Each side can choose a position from 216481 options, how many unique arrangements are possible on the board

Also make it so that the program generates all joint unique positions of white and black arrangements, for example `PPPPPPPP/RRNNBBQK-ppppppqn/rrnbbkpp`

If there is an argument `full`, then print in the `full` folder all possible all joint unique positions of white and black arrangements, for example `PPPPPPPP/RRNNBBQK-ppppppqn/rrnbbkpp`, and in one file there should be no more than 1 million lines. If there is no argument, then print only for white and black in separate files