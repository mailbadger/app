import React, { Fragment, useState, useEffect } from "react";
import qs from "qs";
import { NavLink } from "react-router-dom";
import { Paragraph } from "grommet";
import { mainInstance as axios } from "./axios";

const Unsubscribe = () => {
  const [data, setData] = useState({ message: "" });
  const parsed = qs.parse(window.location.search.slice(1));

  useEffect(() => {
    const callApi = async () => {
      try {
        const res = await axios.post(`/api/unsubscribe`);
        setData(res.data);
      } catch (error) {
        setData(error.response.data);
      }
    };

    callApi();
  }, []);

  if (data.message === "") {
    return <div>Loading...</div>;
  }

  return (
    <Fragment>
      <Paragraph>{data.message}</Paragraph>
      <NavLink to="/dashboard">Go to app &gt;</NavLink>
    </Fragment>
  );
};

export default Unsubscribe;
