import { Button, Typography } from "@mui/material";
import SaveIcon from "@mui/icons-material/Save";
import React from "react";
import * as PusherPushNotifications from "@pusher/push-notifications-web";

const beams = new PusherPushNotifications.Client({
    instanceId: "6e482588-a9a1-45a9-b786-2d367fc69eef"
});
beams
    .start()
    .then((beamsClient) => {
        console.log(beamsClient);
        // beamsClient.getDeviceId();
    })
    // .then((deviceId) =>
    //     console.log("Successfully registered with Beams. Device ID:", deviceId)
    // )
    // .then(() => beams.getDeviceInterests())
    .catch(console.error);

const home = () => {
    return (
        <>
            <Typography variant="h2"></Typography>
            <Button
                startIcon={<SaveIcon />}
                variant="contained"
                color="primary"
                size="small"
            >
                Test Notification
            </Button>
            ;
        </>
    );
};

export default home;
