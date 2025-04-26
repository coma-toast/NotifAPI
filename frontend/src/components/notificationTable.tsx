import React, { useState, useEffect } from "react";
import Box from "@mui/material/Box";
import Collapse from "@mui/material/Collapse";
import IconButton from "@mui/material/IconButton";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import Typography from "@mui/material/Typography";
import Paper from "@mui/material/Paper";
import RefreshIcon from "@mui/icons-material/Refresh";
import DeleteIcon from "@mui/icons-material/Delete";
import KeyboardArrowDownIcon from "@mui/icons-material/KeyboardArrowDown";
import KeyboardArrowUpIcon from "@mui/icons-material/KeyboardArrowUp";
import { NotificationItem } from "../types/notification";

export declare type RefreshFunc = () => void;
export declare type RemoveFunc = (x: string) => void;

interface NotificationProps {
    notifications: NotificationItem[];
    RefreshFunc: RefreshFunc;
    RemoveFunc: RemoveFunc;
}

const NotificationTable = (props: NotificationProps) => {
    const [sorted, setSorted] = useState(props.notifications);

    useEffect(() => {
        let newNotifications = sort(props.notifications);
        setSorted(newNotifications);
    }, [props.notifications]);

    const sort = (x: any[]) => {
        //sort
        return x;
    };

    const remove = (x: string) => {
        sorted.forEach((element, index) => {
            if (element.pub_id === x) {
                sorted.splice(index, 1);
            }
        });
    };

    return (
        <>
            <div className="">
                <CollapsibleTable
                    rows={sorted}
                    refreshFunc={props.RefreshFunc}
                    removeFunc={remove}
                />
            </div>
        </>
    );
};

export default NotificationTable;

function CollapsibleTable(props: {
    rows: NotificationItem[];
    refreshFunc: RefreshFunc;
    removeFunc: RemoveFunc;
}) {
    return (
        <TableContainer component={Paper}>
            <Table aria-label="collapsible table">
                <TableHead>
                    <TableRow>
                        <TableCell>
                            <IconButton onClick={props.refreshFunc}>
                                <RefreshIcon color="primary" />
                            </IconButton>
                        </TableCell>
                        <TableCell>ID</TableCell>
                        <TableCell align="right">Date</TableCell>
                        <TableCell align="right">Source</TableCell>
                        <TableCell align="right">Interests</TableCell>
                        <TableCell align="right">Title</TableCell>
                        <TableCell align="right">Message</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {props.rows.map((row) => (
                        <Row
                            key={row.pub_id}
                            row={row}
                            removeFunc={props.removeFunc}
                        />
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
}

function Row(props: { row: NotificationItem; removeFunc: RemoveFunc }) {
    const { row } = props;
    const [open, setOpen] = React.useState(false);

    return (
        <React.Fragment>
            <TableRow sx={{ "& > *": { borderBottom: "unset" } }}>
                <TableCell>
                    <IconButton
                        aria-label="expand row"
                        size="small"
                        onClick={() => setOpen(!open)}
                    >
                        {open ? (
                            <KeyboardArrowUpIcon />
                        ) : (
                            <KeyboardArrowDownIcon />
                        )}
                    </IconButton>
                </TableCell>
                <TableCell component="th" scope="row">
                    <IconButton onClick={() => props.removeFunc(row.pub_id)}>
                        <DeleteIcon></DeleteIcon>
                    </IconButton>
                </TableCell>
                <TableCell align="right">{row.date}</TableCell>
                <TableCell align="right">{row.source}</TableCell>
                <TableCell align="right">{row.interests}</TableCell>
                <TableCell align="right">{row.title}</TableCell>
                <TableCell align="right">{row.message}</TableCell>
            </TableRow>
            {/* <TableRow>
                <TableCell
                    style={{ paddingBottom: 0, paddingTop: 0 }}
                    colSpan={6}
                >
                    <Collapse in={open} timeout="auto" unmountOnExit>
                        <Box sx={{ margin: 1 }}>
                            <Typography
                                variant="h6"
                                gutterBottom
                                component="div"
                            >
                                History
                            </Typography>
                            <Table size="small" aria-label="purchases">
                                <TableHead>
                                    <TableRow>
                                        <TableCell>Date</TableCell>
                                        <TableCell>Customer</TableCell>
                                        <TableCell align="right">
                                            Amount
                                        </TableCell>
                                        <TableCell align="right">
                                            Total price ($)
                                        </TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {row.history.map((historyRow) => (
                                        <TableRow key={historyRow.date}>
                                            <TableCell
                                                component="th"
                                                scope="row"
                                            >
                                                {historyRow.date}
                                            </TableCell>
                                            <TableCell>
                                                {historyRow.customerId}
                                            </TableCell>
                                            <TableCell align="right">
                                                {historyRow.amount}
                                            </TableCell>
                                            <TableCell align="right">
                                                {Math.round(
                                                    historyRow.amount *
                                                        row.price *
                                                        100
                                                ) / 100}
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </Box>
                    </Collapse>
                </TableCell>
            </TableRow> */}
        </React.Fragment>
    );
}
