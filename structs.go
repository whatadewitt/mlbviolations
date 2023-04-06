package main

import "time"

type ScheduleData struct {
	Copyright            string  `json:"copyright"`
	TotalItems           int     `json:"totalItems"`
	TotalEvents          int     `json:"totalEvents"`
	TotalGames           int     `json:"totalGames"`
	TotalGamesInProgress int     `json:"totalGamesInProgress"`
	Dates                []Date  `json:"dates"`
	Headers              Headers `json:"headers"`
	Status               int     `json:"status"`
}

type GameData struct {
	Copyright        string   `json:"copyright"`
	GamePk           int      `json:"gamePk"`
	Link             string   `json:"link"`
	MetaData         MetaData `json:"metaData"`
	GameData         Game     `json:"gameData"`
	LiveData         LiveData `json:"liveData"`
	ViolationsTotal  int      `json:"violationsTotal,omitEmpty"`
	LastPlayIdx      int      `json:"lastPlayIdx,omitEmpty"`
	LastPlayEventIdx int      `json:"lastPlayEventIdx,omitEmpty"`
	DelayInSeconds   int      `json:"delayInSeconds,omitEmpty"`
}

type Game struct {
	Meta                GameMeta          `json:"game"`
	Datetime            Datetime          `json:"datetime"`
	Status              GameStatus        `json:"status"`
	Teams               GameTeams         `json:"teams"`
	Players             map[string]Player `json:"players"`
	Venue               Venue             `json:"venue"`
	OfficialVenue       OfficialVenue     `json:"officialVenue"`
	Weather             Weather           `json:"weather"`
	GameInfo            GameInfo          `json:"gameInfo"`
	Review              Review            `json:"review"`
	Flags               Flags             `json:"flags"`
	Alerts              []interface{}     `json:"alerts"`
	ProbablePitchers    ProbablePitchers  `json:"probablePitchers"`
	OfficialScorer      Person            `json:"officialScorer"`
	PrimaryDatacaster   Person            `json:"primaryDatacaster"`
	SecondaryDatacaster Person            `json:"secondaryDatacaster"`
}

type GameTeams struct {
	Away Team `json:"away"`
	Home Team `json:"home"`
}

type OfficialVenue struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}

type ProbablePitchers struct {
	Away Player `json:"away"`
	Home Player `json:"home"`
}

type Weather struct {
	Condition string `json:"condition"`
	Temp      string `json:"temp"`
	Wind      string `json:"wind"`
}

type GameInfo struct {
	Attendance          int       `json:"attendance"`
	FirstPitch          time.Time `json:"firstPitch"`
	GameDurationMinutes int       `json:"gameDurationMinutes"`
}

type GameMeta struct {
	Pk              int    `json:"pk"`
	Type            string `json:"type"`
	DoubleHeader    string `json:"doubleHeader"`
	ID              string `json:"id"`
	GamedayType     string `json:"gamedayType"`
	Tiebreaker      string `json:"tiebreaker"`
	GameNumber      int    `json:"gameNumber"`
	CalendarEventID string `json:"calendarEventID"`
	Season          string `json:"season"`
	SeasonDisplay   string `json:"seasonDisplay"`
}

type Datetime struct {
	DateTime     time.Time `json:"dateTime"`
	OriginalDate string    `json:"originalDate"`
	OfficialDate string    `json:"officialDate"`
	DayNight     string    `json:"dayNight"`
	Time         string    `json:"time"`
	Ampm         string    `json:"ampm"`
}

type LiveData struct {
	Plays     Plays     `json:"plays"`
	Linescore Linescore `json:"linescore"`
	Boxscore  Boxscore  `json:"boxscore"`
	Decisions Decisions `json:"decisions"`
	Leaders   Leaders   `json:"leaders"`
}

type Plays struct {
	AllPlays      []Play          `json:"allPlays"`
	About         About           `json:"about"`
	CurrentPlay   Play            `json:"currentPlay"`
	ScoringPlays  []int           `json:"scoringPlays"`
	PlaysByInning []PlaysByInning `json:"playsByInning"`
}

type PlayEvent struct {
	Count       Count     `json:"count"`
	Index       int       `json:"index"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime,omitempty"`
	IsPitch     bool      `json:"isPitch"`
	Type        string    `json:"type"`
	Player      Player    `json:"player,omitempty"`
	Details     Details   `json:"details,omitempty"`
	PitchData   PitchData `json:"pitchData,omitempty"`
	HitData     HitData   `json:"hitData,omitempty"`
	PlayID      string    `json:"playId,omitempty"`
	PitchNumber int       `json:"pitchNumber,omitempty"`
}

type HitData struct {
	LaunchSpeed   float64     `json:"launchSpeed"`
	LaunchAngle   float64     `json:"launchAngle"`
	TotalDistance float64     `json:"totalDistance"`
	Trajectory    string      `json:"trajectory"`
	Hardness      string      `json:"hardness"`
	Location      string      `json:"location"`
	Coordinates   Coordinates `json:"coordinates"`
}

type Coordinates struct {
	AY   float64 `json:"aY"`
	AZ   float64 `json:"aZ"`
	PfxX float64 `json:"pfxX"`
	PfxZ float64 `json:"pfxZ"`
	PX   float64 `json:"pX"`
	PZ   float64 `json:"pZ"`
	VX0  float64 `json:"vX0"`
	VY0  float64 `json:"vY0"`
	VZ0  float64 `json:"vZ0"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	X0   float64 `json:"x0"`
	Y0   float64 `json:"y0"`
	Z0   float64 `json:"z0"`
	AX   float64 `json:"aX"`
}
type Breaks struct {
	BreakAngle    float64 `json:"breakAngle"`
	BreakLength   float64 `json:"breakLength"`
	BreakY        float64 `json:"breakY"`
	SpinRate      int     `json:"spinRate"`
	SpinDirection int     `json:"spinDirection"`
}
type PitchData struct {
	StartSpeed       float64     `json:"startSpeed"`
	EndSpeed         float64     `json:"endSpeed"`
	StrikeZoneTop    float64     `json:"strikeZoneTop"`
	StrikeZoneBottom float64     `json:"strikeZoneBottom"`
	Coordinates      Coordinates `json:"coordinates"`
	Breaks           Breaks      `json:"breaks"`
	Zone             int         `json:"zone"`
	TypeConfidence   float64     `json:"typeConfidence"`
	PlateTime        float64     `json:"plateTime"`
	Extension        float64     `json:"extension"`
}

type Play struct {
	Result      Result      `json:"result"`
	About       About       `json:"about"`
	Count       Count       `json:"count"`
	Matchup     Matchup     `json:"matchup,omitempty"`
	PitchIndex  []int       `json:"pitchIndex"`
	ActionIndex []int       `json:"actionIndex"`
	RunnerIndex []int       `json:"runnerIndex"`
	Runners     []Runners   `json:"runners"`
	PlayEvents  []PlayEvent `json:"playEvents"`
	PlayEndTime time.Time   `json:"playEndTime"`
	AtBatIndex  int         `json:"atBatIndex"`
}

type PlaysByInning struct {
	StartIndex int   `json:"startIndex"`
	EndIndex   int   `json:"endIndex"`
	Top        []int `json:"top"`
	Bottom     []int `json:"bottom"`
	Hits       Hits  `json:"hits"`
}

type Hits struct {
	Away []Matchup `json:"away"`
	Home []Matchup `json:"home"`
}

type Boxscore struct {
	Teams         Teams         `json:"teams"`
	Officials     []Officials   `json:"officials"`
	Info          []Info        `json:"info"`
	PitchingNotes []interface{} `json:"pitchingNotes"`
}

type Info struct {
	Label string `json:"label"`
	Value string `json:"value,omitempty"`
}

type Person struct {
	ID       int    `json:"id"`
	FullName string `json:"fullName"`
	Link     string `json:"link"`
}

type Officials struct {
	Official     Person `json:"official"`
	OfficialType string `json:"officialType"`
}

type MetaData struct {
	Wait          int      `json:"wait"`
	TimeStamp     string   `json:"timeStamp"`
	GameEvents    []string `json:"gameEvents"`
	LogicalEvents []string `json:"logicalEvents"`
}

type Date struct {
	Date                 string         `json:"date"`
	TotalItems           int            `json:"totalItems"`
	TotalEvents          int            `json:"totalEvents"`
	TotalGames           int            `json:"totalGames"`
	TotalGamesInProgress int            `json:"totalGamesInProgress"`
	Games                []ScheduleGame `json:"games"`
	Events               []interface{}  `json:"events"`
}

type ScheduleGame struct {
	SortIndex1             int          `json:"sortIndex1"`
	SortIndex2             int          `json:"sortIndex2"`
	SortIndex3             int          `json:"sortIndex3"`
	SortIndex4             int          `json:"sortIndex4"`
	SortIndex5             int          `json:"sortIndex5"`
	Teams                  Teams        `json:"teams"`
	GamePk                 int          `json:"gamePk"`
	Link                   string       `json:"link"`
	GameType               string       `json:"gameType"`
	Season                 string       `json:"season"`
	GameDate               time.Time    `json:"gameDate"`
	OfficialDate           string       `json:"officialDate"`
	Status                 GameStatus   `json:"status,omitempty"`
	Linescore              Linescore    `json:"linescore,omitempty"`
	Decisions              Decisions    `json:"decisions,omitempty"`
	Venue                  Venue        `json:"venue"`
	Broadcasts             []Broadcasts `json:"broadcasts"`
	Content                Content      `json:"content,omitempty"`
	SeriesStatus           SeriesStatus `json:"seriesStatus,omitempty"`
	IsTie                  bool         `json:"isTie,omitempty"`
	XrefIds                []XrefIds    `json:"xrefIds"`
	GameNumber             int          `json:"gameNumber"`
	PublicFacing           bool         `json:"publicFacing"`
	DoubleHeader           string       `json:"doubleHeader"`
	GamedayType            string       `json:"gamedayType"`
	Tiebreaker             string       `json:"tiebreaker"`
	CalendarEventID        string       `json:"calendarEventID"`
	SeasonDisplay          string       `json:"seasonDisplay"`
	DayNight               string       `json:"dayNight"`
	ScheduledInnings       int          `json:"scheduledInnings"`
	ReverseHomeAwayStatus  bool         `json:"reverseHomeAwayStatus"`
	InningBreakLength      int          `json:"inningBreakLength"`
	GamesInSeries          int          `json:"gamesInSeries"`
	SeriesGameNumber       int          `json:"seriesGameNumber"`
	SeriesDescription      string       `json:"seriesDescription"`
	Review                 Review       `json:"review,omitempty"`
	Flags                  Flags        `json:"flags,omitempty"`
	HomeRuns               []HomeRuns   `json:"homeRuns,omitempty"`
	RecordSource           string       `json:"recordSource"`
	IfNecessary            string       `json:"ifNecessary"`
	IfNecessaryDescription string       `json:"ifNecessaryDescription"`
	GameUtils              GameUtils    `json:"gameUtils"`
	Description            string       `json:"description,omitempty"`
	RescheduleDate         time.Time    `json:"rescheduleDate,omitempty"`
	RescheduleGameDate     string       `json:"rescheduleGameDate,omitempty"`
	Tickets                []Tickets    `json:"tickets,omitempty"`
}

type Teams struct {
	Home GameTeam `json:"home"`
	Away GameTeam `json:"away"`
}

type GameTeam struct {
	IsContextTeam   bool                      `json:"isContextTeam"`
	IsFollowed      bool                      `json:"isFollowed"`
	IsFavorite      bool                      `json:"isFavorite"`
	LeagueRecord    LeagueRecord              `json:"leagueRecord"`
	Score           int                       `json:"score"`
	Team            Team                      `json:"team"`
	IsWinner        bool                      `json:"isWinner"`
	ProbablePitcher Player                    `json:"probablePitcher"`
	SplitSquad      bool                      `json:"splitSquad"`
	SeriesNumber    int                       `json:"seriesNumber"`
	SpringLeague    SpringLeague              `json:"springLeague"`
	Runs            int                       `json:"runs"`
	Hits            int                       `json:"hits"`
	Errors          int                       `json:"errors"`
	LeftOnBase      int                       `json:"leftOnBase"`
	Players         map[string]BoxscorePlayer `json:"players"`
	TeamStats       StatsGroup                `json:"teamStats"`
}

type LeagueRecord struct {
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Pct    string `json:"pct"`
}

type Player struct {
	ID                 int           `json:"id"`
	FullName           string        `json:"fullName"`
	Link               string        `json:"link"`
	FirstName          string        `json:"firstName"`
	LastName           string        `json:"lastName"`
	PrimaryNumber      string        `json:"primaryNumber"`
	BirthDate          string        `json:"birthDate"`
	CurrentAge         int           `json:"currentAge"`
	BirthCity          string        `json:"birthCity"`
	BirthStateProvince string        `json:"birthStateProvince"`
	BirthCountry       string        `json:"birthCountry"`
	Height             string        `json:"height"`
	Weight             int           `json:"weight"`
	Active             bool          `json:"active"`
	AlternateCaptain   bool          `json:"alternateCaptain"`
	Captain            bool          `json:"captain"`
	Rookie             bool          `json:"rookie"`
	PrimaryPosition    Position      `json:"primaryPosition"`
	UseName            string        `json:"useName"`
	MiddleName         string        `json:"middleName"`
	BoxscoreName       string        `json:"boxscoreName"`
	NickName           string        `json:"nickName"`
	Gender             string        `json:"gender"`
	IsPlayer           bool          `json:"isPlayer"`
	IsVerified         bool          `json:"isVerified"`
	DraftYear          int           `json:"draftYear"`
	Stats              []PlayerStats `json:"stats"`
	MlbDebutDate       string        `json:"mlbDebutDate"`
	BatSide            BatSide       `json:"batSide"`
	PitchHand          PitchHand     `json:"pitchHand"`
	NameFirstLast      string        `json:"nameFirstLast"`
	NameSlug           string        `json:"nameSlug"`
	FirstLastName      string        `json:"firstLastName"`
	LastFirstName      string        `json:"lastFirstName"`
	LastInitName       string        `json:"lastInitName"`
	InitLastName       string        `json:"initLastName"`
	FullFMLName        string        `json:"fullFMLName"`
	FullLFMName        string        `json:"fullLFMName"`
	StrikeZoneTop      float64       `json:"strikeZoneTop"`
	StrikeZoneBottom   float64       `json:"strikeZoneBottom"`
}

type BoxscorePlayer struct {
	Person       Person           `json:"person"`
	JerseyNumber string           `json:"jerseyNumber"`
	Position     Position         `json:"position"`
	Status       PlayerGameStatus `json:"status"`
	ParentTeamID int              `json:"parentTeamId"`
	BattingOrder string           `json:"battingOrder"`
	Stats        StatsGroup       `json:"stats"`
	SeasonStats  StatsGroup       `json:"seasonStats"`
	GameStatus   GameStatus       `json:"gameStatus"`
	AllPositions []Position       `json:"allPositions"`
}

type StatsGroup struct {
	Batting  *Stats `json:"batting,omitempty"`
	Pitching *Stats `json:"pitching,omitempty"`
	Fielding *Stats `json:"fielding,omitempty"`
}

type Position struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type Stats struct {
	GamesPlayed            int           `json:"gamesPlayed,omitempty"`
	FlyOuts                int           `json:"flyOuts,omitempty"`
	GroundOuts             int           `json:"groundOuts,omitempty"`
	Runs                   int           `json:"runs,omitempty"`
	Doubles                int           `json:"doubles,omitempty"`
	Triples                int           `json:"triples,omitempty"`
	HomeRuns               int           `json:"homeRuns,omitempty"`
	StrikeOuts             int           `json:"strikeOuts,omitempty"`
	BaseOnBalls            int           `json:"baseOnBalls,omitempty"`
	IntentionalWalks       int           `json:"intentionalWalks,omitempty"`
	Hits                   int           `json:"hits,omitempty"`
	HitByPitch             int           `json:"hitByPitch,omitempty"`
	Avg                    string        `json:"avg,omitempty"`
	AtBats                 int           `json:"atBats,omitempty"`
	Obp                    string        `json:"obp,omitempty"`
	Slg                    string        `json:"slg,omitempty"`
	Ops                    string        `json:"ops,omitempty"`
	CaughtStealing         int           `json:"caughtStealing,omitempty"`
	StolenBases            int           `json:"stolenBases,omitempty"`
	StolenBasePercentage   string        `json:"stolenBasePercentage,omitempty"`
	GroundIntoDoublePlay   int           `json:"groundIntoDoublePlay,omitempty"`
	GroundIntoTriplePlay   int           `json:"groundIntoTriplePlay,omitempty"`
	PlateAppearances       int           `json:"plateAppearances,omitempty"`
	TotalBases             int           `json:"totalBases,omitempty"`
	Rbi                    int           `json:"rbi,omitempty"`
	LeftOnBase             int           `json:"leftOnBase,omitempty"`
	SacBunts               int           `json:"sacBunts,omitempty"`
	SacFlies               int           `json:"sacFlies,omitempty"`
	Babip                  string        `json:"babip,omitempty"`
	CatchersInterference   int           `json:"catchersInterference,omitempty"`
	Pickoffs               int           `json:"pickoffs,omitempty"`
	AtBatsPerHomeRun       string        `json:"atBatsPerHomeRun,omitempty"`
	GamesStarted           int           `json:"gamesStarted,omitempty"`
	AirOuts                int           `json:"airOuts,omitempty"`
	NumberOfPitches        int           `json:"numberOfPitches,omitempty"`
	Era                    string        `json:"era,omitempty"`
	InningsPitched         string        `json:"inningsPitched,omitempty"`
	Wins                   int           `json:"wins,omitempty"`
	Losses                 int           `json:"losses,omitempty"`
	Saves                  int           `json:"saves,omitempty"`
	SaveOpportunities      int           `json:"saveOpportunities,omitempty"`
	Holds                  int           `json:"holds,omitempty"`
	BlownSaves             int           `json:"blownSaves,omitempty"`
	EarnedRuns             int           `json:"earnedRuns,omitempty"`
	Whip                   string        `json:"whip,omitempty"`
	BattersFaced           int           `json:"battersFaced,omitempty"`
	Outs                   int           `json:"outs,omitempty"`
	GamesPitched           int           `json:"gamesPitched,omitempty"`
	CompleteGames          int           `json:"completeGames,omitempty"`
	Shutouts               int           `json:"shutouts,omitempty"`
	PitchesThrown          int           `json:"pitchesThrown,omitempty"`
	Balls                  int           `json:"balls,omitempty"`
	Strikes                int           `json:"strikes,omitempty"`
	StrikePercentage       string        `json:"strikePercentage,omitempty"`
	HitBatsmen             int           `json:"hitBatsmen,omitempty"`
	Balks                  int           `json:"balks,omitempty"`
	WildPitches            int           `json:"wildPitches,omitempty"`
	GroundOutsToAirouts    string        `json:"groundOutsToAirouts,omitempty"`
	WinPercentage          string        `json:"winPercentage,omitempty"`
	PitchesPerInning       string        `json:"pitchesPerInning,omitempty"`
	GamesFinished          int           `json:"gamesFinished,omitempty"`
	StrikeoutWalkRatio     string        `json:"strikeoutWalkRatio,omitempty"`
	StrikeoutsPer9Inn      string        `json:"strikeoutsPer9Inn,omitempty"`
	WalksPer9Inn           string        `json:"walksPer9Inn,omitempty"`
	HitsPer9Inn            string        `json:"hitsPer9Inn,omitempty"`
	RunsScoredPer9         string        `json:"runsScoredPer9,omitempty"`
	HomeRunsPer9           string        `json:"homeRunsPer9,omitempty"`
	InheritedRunners       int           `json:"inheritedRunners,omitempty"`
	InheritedRunnersScored int           `json:"inheritedRunnersScored,omitempty"`
	Note                   string        `json:"note,omitempty"`
	Type                   Type          `json:"type,omitempty"`
	Group                  Group         `json:"group,omitempty"`
	Exemptions             []interface{} `json:"exemptions,omitempty"`
	Splits                 []Splits      `json:"splits,omitempty"`
	Assists                int           `json:"assists,omitempty"`
	PutOuts                int           `json:"putOuts,omitempty"`
	Errors                 int           `json:"errors,omitempty"`
	Chances                int           `json:"chances,omitempty"`
	Fielding               string        `json:"fielding,omitempty"`
	PassedBall             int           `json:"passedBall,omitempty"`
}

type PlayerGameStatus struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type BatSide struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type PitchHand struct {
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
}

type SpringLeague struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Link         string `json:"link"`
	Abbreviation string `json:"abbreviation"`
}

type Team struct {
	SpringLeague    SpringLeague `json:"springLeague"`
	AllStarStatus   string       `json:"allStarStatus"`
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Link            string       `json:"link"`
	Season          int          `json:"season"`
	Venue           Venue        `json:"venue"`
	SpringVenue     SpringVenue  `json:"springVenue"`
	TeamCode        string       `json:"teamCode"`
	FileCode        string       `json:"fileCode"`
	Abbreviation    string       `json:"abbreviation"`
	TeamName        string       `json:"teamName"`
	LocationName    string       `json:"locationName"`
	FirstYearOfPlay string       `json:"firstYearOfPlay"`
	League          League       `json:"league"`
	Division        Division     `json:"division"`
	Sport           Sport        `json:"sport"`
	ShortName       string       `json:"shortName"`
	FranchiseName   string       `json:"franchiseName"`
	ClubName        string       `json:"clubName"`
	Active          bool         `json:"active"`
}

type SpringVenue struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
}

// TODO: these are pretty generic
// share a struct between them?
type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Division struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Sport struct {
	ID   int    `json:"id"`
	Link string `json:"link"`
	Name string `json:"name"`
}

type GameStatus struct {
	AbstractGameState string `json:"abstractGameState"`
	CodedGameState    string `json:"codedGameState"`
	DetailedState     string `json:"detailedState"`
	StatusCode        string `json:"statusCode"`
	StartTimeTBD      bool   `json:"startTimeTBD"`
	Reason            string `json:"reason"`
	AbstractGameCode  string `json:"abstractGameCode"`
	IsCurrentBatter   bool   `json:"isCurrentBatter"`
	IsCurrentPitcher  bool   `json:"isCurrentPitcher"`
	IsOnBench         bool   `json:"isOnBench"`
	IsSubstitute      bool   `json:"isSubsitute"`
}

type Linescore struct {
	CurrentInning        int       `json:"currentInning"`
	CurrentInningOrdinal string    `json:"currentInningOrdinal"`
	InningState          string    `json:"inningState"`
	InningHalf           string    `json:"inningHalf"`
	IsTopInning          bool      `json:"isTopInning"`
	ScheduledInnings     int       `json:"scheduledInnings"`
	Innings              []Innings `json:"innings"`
	Teams                Teams     `json:"teams"`
	Defense              Defense   `json:"defense"`
	Offense              Offense   `json:"offense"`
	Balls                int       `json:"balls"`
	Strikes              int       `json:"strikes"`
	Outs                 int       `json:"outs"`
}

type Innings struct {
	Num        int         `json:"num"`
	OrdinalNum string      `json:"ordinalNum"`
	Home       InningStats `json:"home"`
	Away       InningStats `json:"away"`
}

type InningStats struct {
	Runs       int `json:"runs"`
	Hits       int `json:"hits"`
	Errors     int `json:"errors"`
	LeftOnBase int `json:"leftOnBase"`
}

type Offense struct {
	Batter       Player `json:"batter"`
	OnDeck       Player `json:"onDeck"`
	InHole       Player `json:"inHole"`
	BattingOrder int    `json:"battingOrder"`
	Team         Team   `json:"team"`
}

type Defense struct {
	Pitcher      Player `json:"pitcher"`
	Catcher      Player `json:"catcher"`
	First        Player `json:"first"`
	Second       Player `json:"second"`
	Third        Player `json:"third"`
	Shortstop    Player `json:"shortstop"`
	Left         Player `json:"left"`
	Center       Player `json:"center"`
	Right        Player `json:"right"`
	BattingOrder int    `json:"battingOrder"`
}

type Decisions struct {
	Winner *Player `json:"winner,omitempty"`
	Loser  *Player `json:"loser,omitempty"`
	Save   *Player `json:"save,omitempty"`
}

type Venue struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Link     string   `json:"link"`
	Location Location `json:"location"`
	Active   bool     `json:"active"`
}

type Location struct {
	Address1           string             `json:"address1"`
	City               string             `json:"city"`
	State              string             `json:"state"`
	StateAbbrev        string             `json:"stateAbbrev"`
	PostalCode         string             `json:"postalCode"`
	DefaultCoordinates DefaultCoordinates `json:"defaultCoordinates"`
	Country            string             `json:"country"`
	Phone              string             `json:"phone"`
}

type DefaultCoordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Broadcasts struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Type            string          `json:"type"`
	Language        string          `json:"language"`
	HomeAway        string          `json:"homeAway"`
	CallSign        string          `json:"callSign"`
	IsNational      bool            `json:"isNational,omitempty"`
	VideoResolution VideoResolution `json:"videoResolution,omitempty"`
}

type VideoResolution struct {
	Code            string `json:"code"`
	ResolutionShort string `json:"resolutionShort"`
	ResolutionFull  string `json:"resolutionFull"`
}

type Content struct {
	Link       string     `json:"link"`
	Editorial  Editorial  `json:"editorial"`
	Media      Media      `json:"media"`
	Highlights Highlights `json:"highlights"`
	Summary    Summary    `json:"summary"`
	GameNotes  GameNotes  `json:"gameNotes"`
}

type Editorial struct {
}

type Highlights struct {
}

type Summary struct {
	HasPreviewArticle  bool `json:"hasPreviewArticle"`
	HasRecapArticle    bool `json:"hasRecapArticle"`
	HasWrapArticle     bool `json:"hasWrapArticle"`
	HasHighlightsVideo bool `json:"hasHighlightsVideo"`
}

type GameNotes struct {
}

type SeriesStatus struct {
	GameNumber       int      `json:"gameNumber"`
	TotalGames       int      `json:"totalGames"`
	IsTied           bool     `json:"isTied"`
	IsOver           bool     `json:"isOver"`
	Wins             int      `json:"wins"`
	Losses           int      `json:"losses"`
	WinningTeam      GameTeam `json:"winningTeam"`
	LosingTeam       GameTeam `json:"losingTeam"`
	Description      string   `json:"description"`
	ShortDescription string   `json:"shortDescription"`
	Result           string   `json:"result"`
	ShortName        string   `json:"shortName"`
}

type XrefIds struct {
	XrefID   string `json:"xrefId"`
	XrefType string `json:"xrefType"`
}

type Review struct {
	HasChallenges bool       `json:"hasChallenges"`
	Away          TeamReview `json:"away"`
	Home          TeamReview `json:"home"`
}

type TeamReview struct {
	Used      int `json:"used"`
	Remaining int `json:"remaining"`
}

type Flags struct {
	NoHitter            bool `json:"noHitter"`
	PerfectGame         bool `json:"perfectGame"`
	AwayTeamNoHitter    bool `json:"awayTeamNoHitter"`
	AwayTeamPerfectGame bool `json:"awayTeamPerfectGame"`
	HomeTeamNoHitter    bool `json:"homeTeamNoHitter"`
	HomeTeamPerfectGame bool `json:"homeTeamPerfectGame"`
}

type HomeRuns struct {
	Result      Result        `json:"result"`
	About       About         `json:"about"`
	Count       Count         `json:"count"`
	Matchup     Matchup       `json:"matchup"`
	PitchIndex  []interface{} `json:"pitchIndex"`
	ActionIndex []interface{} `json:"actionIndex"`
	RunnerIndex []interface{} `json:"runnerIndex"`
	Runners     []interface{} `json:"runners"`
	PlayEvents  []interface{} `json:"playEvents"`
}

type Matchup struct {
	Team                Team          `json:"team"`
	InningStats         int           `json:"inning"`
	Batter              Player        `json:"batter"`
	Pitcher             Player        `json:"pitcher"`
	BatSide             BatSide       `json:"batSide"`
	PitchHand           PitchHand     `json:"pitchHand"`
	BatterHotColdZones  []interface{} `json:"batterHotColdZones"`
	PitcherHotColdZones []interface{} `json:"pitcherHotColdZones"`
	Splits              Splits        `json:"splits"`
	Coordinates         Point         `json:"coordinates"`
	Type                string        `json:"type"`
	Description         string        `json:"description"`
}

type Zone struct {
	Zone  string `json:"zone"`
	Color string `json:"color"`
	Temp  string `json:"temp"`
	Value string `json:"value"`
}

type Stat struct {
	Name  string `json:"name"`
	Zones []Zone `json:"zones"`
}

type Splits struct {
	Stat Stat `json:"stat"`
}

type Point struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type GameUtils struct {
	IsAllStar                   bool `json:"isAllStar"`
	IsCancelled                 bool `json:"isCancelled"`
	IsClassicDoubleHeader       bool `json:"isClassicDoubleHeader"`
	IsCompletedEarly            bool `json:"isCompletedEarly"`
	IsDelayed                   bool `json:"isDelayed"`
	IsDoubleHeader              bool `json:"isDoubleHeader"`
	IsNonDoubleHeaderTBD        bool `json:"isNonDoubleHeaderTBD"`
	IsExhibition                bool `json:"isExhibition"`
	IsFinal                     bool `json:"isFinal"`
	IsForfeit                   bool `json:"isForfeit"`
	IsFreeGame                  bool `json:"isFreeGame"`
	IsInstantReplay             bool `json:"isInstantReplay"`
	IsLive                      bool `json:"isLive"`
	IsManagerChallenge          bool `json:"isManagerChallenge"`
	IsNoHitter                  bool `json:"isNoHitter"`
	IsPerfectGame               bool `json:"isPerfectGame"`
	IsPostponed                 bool `json:"isPostponed"`
	IsPreview                   bool `json:"isPreview"`
	IsSplitTicketDoubleHeader   bool `json:"isSplitTicketDoubleHeader"`
	IsSpring                    bool `json:"isSpring"`
	IsSuspended                 bool `json:"isSuspended"`
	IsSuspendedOnDate           bool `json:"isSuspendedOnDate"`
	IsSuspendedResumptionOnDate bool `json:"isSuspendedResumptionOnDate"`
	IsTBD                       bool `json:"isTBD"`
	IsUmpireReview              bool `json:"isUmpireReview"`
	IsWarmup                    bool `json:"isWarmup"`
	IsPostSeason                bool `json:"isPostSeason"`
	IsTieBreaker                bool `json:"isTieBreaker"`
	IsPostSeasonReady           bool `json:"isPostSeasonReady"`
	IsWildCard                  bool `json:"isWildCard"`
	IsDivisionSeries            bool `json:"isDivisionSeries"`
	IsChampionshipSeries        bool `json:"isChampionshipSeries"`
	IsWorldSeries               bool `json:"isWorldSeries"`
	IsPreGameDelay              bool `json:"isPreGameDelay"`
	IsInGameDelay               bool `json:"isInGameDelay"`
	HasContextTeam              bool `json:"hasContextTeam"`
	HasFavorites                bool `json:"hasFavorites"`
	HasMostFavorite             bool `json:"hasMostFavorite"`
	HasFollowed                 bool `json:"hasFollowed"`
}

type Tickets struct {
	TicketType  string      `json:"ticketType"`
	TicketLinks TicketLinks `json:"ticketLinks"`
}

type TicketLinks struct {
	Home string `json:"home"`
}

type Type struct {
	DisplayName string `json:"displayName"`
}
type Group struct {
	DisplayName string `json:"displayName"`
}

type PlayerStats struct {
	Type       Type          `json:"type"`
	Group      Group         `json:"group"`
	Exemptions []interface{} `json:"exemptions"`
	Stats      Stats         `json:"stats,omitempty"`
}

type Items struct {
	ID               int    `json:"id"`
	ContentID        string `json:"contentId"`
	MediaID          string `json:"mediaId"`
	MediaState       string `json:"mediaState"`
	MediaFeedType    string `json:"mediaFeedType"`
	MediaFeedSubType string `json:"mediaFeedSubType"`
	CallLetters      string `json:"callLetters"`
	FoxAuthRequired  bool   `json:"foxAuthRequired"`
	TbsAuthRequired  bool   `json:"tbsAuthRequired"`
	EspnAuthRequired bool   `json:"espnAuthRequired"`
	Fs1AuthRequired  bool   `json:"fs1AuthRequired"`
	MlbnAuthRequired bool   `json:"mlbnAuthRequired"`
	FreeGame         bool   `json:"freeGame"`
	GameDate         string `json:"gameDate"`
}

type Epg struct {
	Title string  `json:"title"`
	Items []Items `json:"items"`
}

type FeaturedMedia struct {
	ID string `json:"id"`
}

type Media struct {
	Epg           []Epg         `json:"epg"`
	FeaturedMedia FeaturedMedia `json:"featuredMedia"`
	FreeGame      bool          `json:"freeGame"`
	EnhancedGame  bool          `json:"enhancedGame"`
}

type Result struct {
	Type        string `json:"type"`
	Event       string `json:"event"`
	Description string `json:"description"`
	Rbi         int    `json:"rbi"`
	AwayScore   int    `json:"awayScore"`
	HomeScore   int    `json:"homeScore"`
}
type About struct {
	AtBatIndex       int       `json:"atBatIndex"`
	HalfInning       string    `json:"halfInning"`
	IsTopInning      bool      `json:"isTopInning"`
	Inning           int       `json:"inning"`
	StartTime        time.Time `json:"startTime"`
	EndTime          time.Time `json:"endTime"`
	IsComplete       bool      `json:"isComplete"`
	IsScoringPlay    bool      `json:"isScoringPlay"`
	HasReview        bool      `json:"hasReview"`
	HasOut           bool      `json:"hasOut"`
	CaptivatingIndex int       `json:"captivatingIndex"`
}

type Count struct {
	Balls   int `json:"balls"`
	Strikes int `json:"strikes"`
	Outs    int `json:"outs"`
}

type Headers struct {
	XSTITCHSHA   string `json:"X-STITCH-SHA"`
	XSTITCHCACHE string `json:"X-STITCH-CACHE"`
	CacheControl string `json:"Cache-Control"`
	ContentType  string `json:"Content-Type"`
}

type Leaders struct {
	HitDistance interface{} `json:"hitDistance"`
	HitSpeed    interface{} `json:"hitSpeed"`
	PitchSpeed  interface{} `json:"pitchSpeed"`
}

type Details struct {
	Description   string `json:"description"`
	Event         string `json:"event"`
	EventType     string `json:"eventType"`
	AwayScore     int    `json:"awayScore"`
	HomeScore     int    `json:"homeScore"`
	IsScoringPlay bool   `json:"isScoringPlay"`
	HasReview     bool   `json:"hasReview"`
	Call          Call   `json:"call"`
}

type Call struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

type Credits struct {
	Player   Player   `json:"player"`
	Position Position `json:"position"`
	Credit   string   `json:"credit"`
}

type Movement struct {
	OriginBase interface{} `json:"originBase"`
	Start      interface{} `json:"start"`
	End        interface{} `json:"end"`
	OutBase    string      `json:"outBase"`
	IsOut      bool        `json:"isOut"`
	OutNumber  int         `json:"outNumber"`
}

type Runners struct {
	Movement Movement  `json:"movement"`
	Details  Details   `json:"details"`
	Credits  []Credits `json:"credits"`
}

type GameCache struct {
	Status      string   `json:"status"`
	PitchedKeys []string `json:"pitched"`
}

type CloserNews struct {
	Tweet          string
	Game           ScheduleGame
	DelayInSeconds int
}

type ViolationNotification struct {
	Tweet string
}
