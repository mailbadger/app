import { useReducer, useEffect, useState, useRef } from "react";
import { mainInstance as axios } from "../axios";

const dataFetchReducer = (state, action) => {
  switch (action.type) {
    case "REQUEST_INIT":
      return {
        ...state,
        isLoading: true,
        isError: false,
      };
    case "REQUEST_SUCCESS":
      return {
        ...state,
        isLoading: false,
        isError: false,
        data: action.payload,
      };
    case "REQUEST_FAILURE":
      return {
        ...state,
        isLoading: false,
        isError: true,
        data: action.payload,
      };
    default:
      throw new Error();
  }
};

const defaultOpts = {};

const useDataApi = (initialOpts = defaultOpts, initialData, skipFirst = false) => {
  const firstUpdate = useRef(skipFirst);
  const [opts, setOpts] = useState(initialOpts);

  const [state, dispatch] = useReducer(dataFetchReducer, {
    isLoading: false,
    isError: false,
    data: initialData,
  });

  useEffect(() => {
    if (firstUpdate.current) {
      firstUpdate.current = false;
      return;
    }

    let didCancel = false;

    const fetchData = async () => {
      dispatch({ type: "REQUEST_INIT" });

      try {
        const result = await axios(opts);

        if (!didCancel) {
          dispatch({ type: "REQUEST_SUCCESS", payload: result.data });
        }
      } catch (error) {
        if (!didCancel) {
          dispatch({ type: "REQUEST_FAILURE", error, payload: error.response.data });
        }
      }
    };

    fetchData();

    return () => {
      didCancel = true;
    };
  }, [opts]);

  const callApi = (opts) => {
    setOpts(opts);
  };

  return [state, callApi];
};

export default useDataApi;
