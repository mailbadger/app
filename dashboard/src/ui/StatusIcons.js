import React from "react";
import {
  StatusGood,
  StatusWarning,
  StatusCritical,
  StatusInfo,
} from "grommet-icons";

const StatusIcon = {
  "status-ok": <StatusGood />,
  "status-warning": <StatusWarning />,
  "status-error": <StatusCritical />,
  "status-critical": <StatusCritical />,
  "status-info": <StatusInfo />,
};

export default StatusIcon;
