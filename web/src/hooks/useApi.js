import { useReducer, useEffect, useState } from "react";
import axios from "axios";

const dataFetchReducer = (state, action) => {
  switch (action.type) {
    case "REQUEST_INIT":
      return {
        ...state,
        isLoading: true,
        isError: false
      };
    case "REQUEST_SUCCESS":
      return {
        ...state,
        isLoading: false,
        isError: false,
        data: action.payload
      };
    case "REQUEST_FAILURE":
      return {
        ...state,
        isLoading: false,
        isError: true
      };
    default:
      throw new Error();
  }
};

const defaultOpts = {};

const useDataApi = (initialOpts = defaultOpts, initialData) => {
  const [opts, setOpts] = useState(initialOpts);

  const [state, dispatch] = useReducer(dataFetchReducer, {
    isLoading: false,
    isError: false,
    data: initialData
  });

  const { readBody = body => body.json(), ...input } = opts;

  useEffect(() => {
    let didCancel = false;

    const fetchData = async () => {
      dispatch({ type: "REQUEST_INIT" });

      try {
        const result = await axios(input);
        const payload = await readBody(result);

        if (!didCancel) {
          dispatch({ type: "REQUEST_SUCCESS", payload: payload });
        }
      } catch (error) {
        if (!didCancel) {
          dispatch({ type: "REQUEST_FAILURE" });
        }
      }
    };

    fetchData();

    return () => {
      didCancel = true;
    };
  }, [opts]);

  const callApi = opts => {
    setOpts(opts);
  };

  return { ...state, callApi };
};

export default useDataApi;
