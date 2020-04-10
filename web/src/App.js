import React, { useState, useEffect } from "react";
import Intro from "./intro";
import Game from "./game";
import "./App.css";

function useGameState() {
    const gameInitialState = () => window.localStorage.getItem("game") || null;
    const [game, setGame] = useState(gameInitialState);
    useEffect(() => {
        if (game === null) {
            window.localStorage.removeItem("game");
        } else {
            window.localStorage.setItem("game", game);
        }
    }, [game]);

    const playerInitialState = () =>
        window.localStorage.getItem("player") || null;
    const [player, setPlayer] = useState(playerInitialState);
    useEffect(() => {
        if (player === null) {
            window.localStorage.removeItem("player");
        } else {
            window.localStorage.setItem("player", player);
        }
    }, [player]);

    return [game, player, setGame, setPlayer];
}

function App() {
    const [game, player, setGame, setPlayer] = useGameState();
    const [started, setStarted] = useState(false);
    const [timer, setTimer] = useState(null);
    const [gameState, setGameState] = useState({});

    useEffect(() => {
        if (!player || !game) {
            return;
        }
        if (timer !== null) {
            return;
        }
        const t = setInterval(() => {
            if (!game || !player) {
                return;
            }
            fetch(`http://localhost:8000/game/${game}/${player}`)
                .then((res) => {
                    if (!res.ok) {
                        throw res.statusText; // TODO: manage the error
                    }
                    return res.json();
                })
                .then((res) => {
                    setGameState(res);
                    if (res.game_started === true) {
                        setStarted(true);
                        clearInterval(t);
                    }
                })
                .catch((err) => console.error(err));
        }, 3000);
        setTimer(t);
        return () => {
            if (timer !== null) {
                console.log("cleanup");
            }
        };
    }, [game, player, started, timer]);

    return (
        <div className="App">
            <Debug
                game={game}
                player={player}
                started={started}
                gameState={gameState}
            />
            {game && player && started ? (
                <Game gameState={gameState} />
            ) : (
                <Intro
                    setGame={setGame}
                    game={game}
                    setPlayer={setPlayer}
                    player={player}
                    waitingPlayers={timer !== null}
                    gameState={gameState}
                />
            )}
        </div>
    );
}

function Debug({ game, player, started, gameState }) {
    const clear = () => {
        ["player", "game"].forEach((key) => {
            console.log(key);
            window.localStorage.removeItem(key);
            location.reload(); //eslint-disable-line
        });
    };
    return (
        <div className="debug">
            <div>Current game: {game}</div>
            <div>Current player: {player}</div>
            <div>Game started: {started ? "true" : "false"}</div>
            <div>Prize cards: {gameState.prize_cards}</div>
            {game == null || player === null ? (
                <span />
            ) : (
                <div>
                    <a
                        href={`http://localhost:8000/game/${game}/${player}`}
                        target="_blank"
                        rel="noopener noreferrer"
                    >
                        Debug
                    </a>
                    <br />
                    <button onClick={() => clear()}>Clear</button>
                </div>
            )}
        </div>
    );
}

export default App;
