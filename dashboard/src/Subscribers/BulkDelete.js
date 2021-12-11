import React, { useEffect, useContext } from "react";
import { Box, Heading, Markdown } from "grommet";
import Uppy from "@uppy/core";
import AwsS3 from "@uppy/aws-s3";
import { DragDrop, StatusBar } from "@uppy/react";
import qs from "qs";

import "@uppy/core/dist/style.css";
import "@uppy/drag-drop/dist/style.css";
import "@uppy/status-bar/dist/style.css";

import { mainInstance as axios } from "../axios";
import { NotificationsContext } from "../Notifications/context";
import { endpoints } from "../network/endpoints";

const Content = `
CSV format:

- Columns should be separated by comma
- Include the **Email** header in the file

**Email** |
--- |
john@example.com |
jane@example.com |
`;

const ImportSubscribers = () => {
  const { createNotification } = useContext(NotificationsContext);

  const uppy = Uppy({
    restrictions: {
      maxNumberOfFiles: 1,
      allowedFileTypes: ["text/csv"],
    },
  });
  uppy.use(AwsS3, {
    async getUploadParameters(file) {
      try {
        const res = await axios.post(
          endpoints.signInS3,
          qs.stringify({
            filename: file.name,
            contentType: file.type,
            action: "remove",
          })
        );

        return res.data;
      } catch (error) {
        let msg = "Unable to upload file. Please try again.";
        if (error.response) {
          msg = error.response.data.message;
        }

        createNotification(msg, "status-error");
      }
    },
  });

  uppy.on("upload-success", async (file) => {
    try {
      const res = await axios.post(
        endpoints.deleteSubscribersBulk,
        qs.stringify(
          {
            filename: file.name,
          },
          { arrayFormat: "brackets" }
        )
      );

      createNotification(res.data.message, "status-ok");
    } catch (error) {
      let msg = "Unable to remove subscribers. Please try again.";
      if (error.response) {
        msg = error.response.data.message;
      }

      createNotification(msg, "status-error");
    }

    uppy.reset();
  });

  useEffect(() => {
    return () => {
      uppy.close();
    };
  }, [uppy]);

  return (
    <Box direction="column" margin="medium" animation="fadeIn">
      <Box pad={{ left: "medium" }} margin={{ bottom: "small" }}>
        <Heading level="2">Delete in bulk</Heading>
      </Box>
      <Box round background="white" pad="medium" width="50%" alignSelf="start">
        <Markdown>{Content}</Markdown>
        <Box margin={{ top: "large" }}>
          <DragDrop
            width="100%"
            height="100%"
            note="Only CSV files are allowed"
            uppy={uppy}
          />
          <StatusBar hideAfterFinish={false} showProgressDetails uppy={uppy} />
        </Box>
      </Box>
    </Box>
  );
};

export default ImportSubscribers;
