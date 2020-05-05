import React from "react";
import { StatusGood, StatusWarning, StatusCritical } from "grommet-icons";

const StatusIcon = {
  "status-ok": <StatusGood />,
  "status-warning": <StatusWarning />,
  "status-error": <StatusCritical />,
  "status-critical": <StatusCritical />,
};

export default StatusIcon;
