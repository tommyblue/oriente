import React, { useState } from "react";

import cover from "./assets/cover.png";
import "./intro.css";

function Intro({
    setGame,
    setPlayer,
    player,
    game,
    waitingPlayers,
    gameState,
}) {
    const [showOption, setShowOption] = useState(null);
    return (
        <div>
            <img src={cover} alt="Oriente" className="cover" />
            {player && game && waitingPlayers ? (
                <Waiting gameState={gameState} />
            ) : (
                <div className="buttons">
                    <NewGame
                        setGame={setGame}
                        setPlayer={setPlayer}
                        showOption={showOption}
                        setShowOption={setShowOption}
                    />
                    <JoinGame
                        setGame={setGame}
                        setPlayer={setPlayer}
                        showOption={showOption}
                        setShowOption={setShowOption}
                    />
                </div>
            )}
        </div>
    );
}

function Waiting({ gameState }) {
    const s = gameState.players
        ? `(${gameState.active_players}/${gameState.players.length})`
        : "";
    return <div className="waiting">Waiting...{s}</div>;
}

function JoinGame({ setGame, setPlayer, showOption, setShowOption }) {
    const [gameID, setGameID] = useState(null);
    const joinGame = () => {
        fetch(`http://localhost:8000/game/${gameID}`)
            .then((res) => {
                if (!res.ok) {
                    throw res.statusText; // TODO: manage the error
                }
                return res.json();
            })
            .then((res) => {
                setGame(res.game);
                setPlayer(res.player);
            })
            .catch((err) => console.error(err));
    };
    if (showOption !== null && showOption !== "join") {
        return <span />;
    }
    if (showOption === "join") {
        return (
            <div>
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        joinGame();
                    }}
                >
                    <section>
                        <label>Game ID:</label>
                    </section>
                    <section>
                        <input
                            type="text"
                            onChange={(e) => setGameID(e.target.value)}
                        />
                    </section>
                    <section>
                        <button className="button">Join!</button>
                    </section>
                </form>
                <button className="button" onClick={() => setShowOption(null)}>
                    &laquo; Back
                </button>
            </div>
        );
    }
    return (
        <button className="button" onClick={() => setShowOption("join")}>
            Join game
        </button>
    );
}

function NewGame({ setGame, setPlayer, showOption, setShowOption }) {
    const [nPlayers, setNPlayers] = useState(4);
    const setGameId = () => {
        fetch(`http://localhost:8000/game/new/${nPlayers}`)
            .then((res) => res.json())
            .then((res) => {
                setGame(res.game);
                setPlayer(res.player);
            });
    };
    if (showOption !== null && showOption !== "new") {
        return <span />;
    }

    if (showOption === "new") {
        return (
            <div>
                <form
                    onSubmit={(e) => {
                        e.preventDefault();
                        setGameId();
                    }}
                >
                    <section>
                        <label htmlFor="players-number">
                            Number of players:{" "}
                        </label>
                    </section>
                    <section>
                        <input
                            id="players-number"
                            type="number"
                            min="4"
                            max="12"
                            value={nPlayers}
                            onChange={(e) => setNPlayers(e.target.value)}
                        />
                    </section>
                    <section>
                        <button className="button">Start!</button>
                    </section>
                </form>
                <button className="button" onClick={() => setShowOption(null)}>
                    &laquo; Back
                </button>
            </div>
        );
    }
    return (
        <button className="button" onClick={() => setShowOption("new")}>
            New game
        </button>
    );
}

export default Intro;
