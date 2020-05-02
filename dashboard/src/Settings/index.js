import React from "react";
import AddSesKeys from "./AddSesKeys";
import ChangePassword from "./ChangePassword";
import { CustomTabs } from "../ui";

const tabs = [
  {
    title: "Email Transport",
    children: <AddSesKeys />,
  },
  {
    title: "Account",
    children: <ChangePassword />,
  },
];

const Settings = () => <CustomTabs tabs={tabs} />;

export default Settings;
