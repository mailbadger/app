import React, { Fragment, useState, useEffect } from "react";
import PropTypes from "prop-types";
import { NavLink } from "react-router-dom";
import { Paragraph } from "grommet";
import axios from "axios";

const VerifyEmail = props => {
  const [data, setData] = useState({ message: "" });
  const {
    match: { params }
  } = props;

  useEffect(() => {
    const callApi = async () => {
      try {
        const res = await axios.put(`/api/verify-email/${params.token}`);
        setData(res.data);
      } catch (error) {
        setData(error.response.data);
      }
    };

    callApi();
  });

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

VerifyEmail.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      token: PropTypes.string
    })
  })
};

export default VerifyEmail;
