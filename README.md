## Battlesloop Protocol

We need to be able to send/receive fires, hits, and misses at certain board positions

In general, we use `:` as delimiters and `_` followed by a letter followed by another `_` to indicate data types

### Components Types -- Messages may contain these

#### Board Position
contains a row (letter) and column (integer)
- Row: \[A-J\]
- Column: \[1-10\]

in messages, takes the form: `<row>-<column>:`

#### Ship Type
consists of a single integer
`<int>:`

### Message Type

#### Positional Messages
format -- `_type_position` 
- type is in lowercase letters
example: `_h_A-7:` -- hit on A7

types:
- `h`: hit
- `m`: miss
- `f`: fire

#### Ship Messages
sank -- `a` (for abyss)
- example: `_a_10:` -- sank ship # 10

#### Game State Messages
`g` for game

win, lose, end (no contest, game ending for some other reason)
`_g_win:`
`_g_lose:`
`_g_end:`

#### Connection-Related Messages
`c` for connection

begin, end, heartbeat, your turn, my turn
`_c_begin:`
`_c_end:`
`_c_mturn:`
`_c_yturn:`

