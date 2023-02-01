// const fetchHistory = async (count: string) => {
//         fetch("/api/recent/" + count)
//         .then((response) => response.json())
//         // return response.json
// export function fetchHistory() {}

// function fetchData(endpoint: string) {

// }

import React, { useState, useEffect } from "react";

export async function registerUser(userdata: FormData) {
    // const [data, setData] = useState(null);
    // const [loading, setLoading] = useState(true);
    // const [error, setError] = useState(null);

    // useEffect(() => {
    //     fetch(`/register`)
    //         .then((response) => response.json())
    //         .then((usefulData) => {
    //             console.log(usefulData);
    //             setLoading(false);
    //             setData(usefulData);
    //         })
    //         .catch((e) => {
    //             console.error(`An error occurred: ${e}`);
    //         });
    // }, []);
    try {
        const res = await fetch(`/register`, {
            method: "POST",
            body: userdata
            // headers: {
            // "Content-Type": "application/json"
            // 'Content-Type': 'application/x-www-form-urlencoded',
            // }
        });

        if (!res.ok) {
            const message = `An error has occurred: ${res.status} - ${res.statusText}`;
            throw new Error(message);
        }

        const data = await res.json();
        return data;
    } catch (e) {
        return JSON.stringify("Error registering user" + e);
    }
}
