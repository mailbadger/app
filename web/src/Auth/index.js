import React from "react";
import { Formik } from "formik";
import Form from "./Form";
import { string, object } from "yup";
import axios from "axios";
import qs from "qs";

const loginValidation = object().shape({
  email: string()
    //.email('Please enter a valid email')
    .required("Please enter an email"),
  password: string().required("Please enter your password")
});

const Auth = () => {
  const handleSubmit = (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        const result = await axios.post(
          "/api/authenticate",
          qs.stringify({
            username: values.email,
            password: values.password
          }),
          {
            headers: {
              "Content-Type": "application/x-www-form-urlencoded"
            }
          }
        );

        console.log(result.data);
      } catch (error) {
        setErrors(error.response.data);
      }
    };

    callApi();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  return (
    <Formik
      onSubmit={handleSubmit}
      validationSchema={loginValidation}
      render={props => <Form {...props} />}
    />
  );
};

export default Auth;
