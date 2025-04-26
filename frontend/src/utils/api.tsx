import React, { useState, useEffect } from "react";

export declare type RefreshFunc = () => void;

interface HistoryProps {
    notifications: any[];
    RefreshFunc: RefreshFunc;
}

const HistoryComponent = (props: HistoryProps) => {
    const [sorted, setSorted] = useState(props.notifications);
    useEffect(() => {
        let newNotifications = sort(props.notifications);
        setSorted(newNotifications);
    }, [props.notifications]);

    const sort = (x: any[]) => {
        //sort
        return x;
    };

    // let newNotifications = sort(props.notifications);
    // setSorted(newNotifications);

    const notificationElements = [] as any[];
    sorted.forEach((notification: any) => {
        notificationElements.push(<div>{notification.message}</div>);
    });

    return (
        // pretend table
        <div className="add-post-container">
            <button onClick={props.RefreshFunc}></button>
            <div>{notificationElements}</div>
        </div>
    );
};

export default HistoryComponent;
