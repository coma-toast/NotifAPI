import React, { useState, useEffect } from "react";
import NotificationTable from "../components/notificationTable";

const History = () => {
    const [count, setCount] = useState("10");
    const [notifications, setNotifications] = useState([]);
    const fetchHistory = async (count: string) => {
        fetch("/api/recent/" + count)
            .then((response) => response.json())
            .then((data) => {
                console.log(data);
                setNotifications(data);
            })
            .catch((err) => {
                console.log(err.message);
            });
    };
    useEffect(() => {
        fetchHistory("10");
    }, []);
    const getHistory = () => {
        fetchHistory("10");
    };

    return (
        <div className="history">
            <form onSubmit={getHistory}>
                <input
                    type="text"
                    className="form-control"
                    value={count}
                    onChange={(e) => setCount(e.target.value)}
                />
                <button type="submit">refresh</button>
            </form>

            <NotificationTable
                notifications={notifications}
                RefreshFunc={getHistory}
            />
        </div>
    );
};

export default History;
