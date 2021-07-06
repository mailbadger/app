import React, { useEffect, useContext, useMemo, useState } from "react";
import { Box, Markdown, Text, Select } from "grommet";
import Uppy from "@uppy/core";
import AwsS3 from "@uppy/aws-s3";
import {
  //  DragDrop,
  StatusBar,
  Dashboard,
} from "@uppy/react";
import qs from "qs";

import "@uppy/core/dist/style.css";
import "@uppy/drag-drop/dist/style.css";
import "@uppy/status-bar/dist/style.css";

import { mainInstance as axios } from "../axios";
import { NotificationsContext } from "../Notifications/context";
import styled from "styled-components";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDownload,
  faCloudUploadAlt,
} from "@fortawesome/free-solid-svg-icons";
import { useApi } from "../hooks";

const StyledImportButton = styled(Box)`
  div {
    border-radius: 20px;
    cursor: pointer;
    display: flex;
    background-color: black;
    color: white;
    ${(props) =>
      props.disabled
        ? `
			pointer-events:none;
			opacity: 0.5;
		`
        : ""};
  }
`;
const StyledDragDrop = styled(Box)`

button {
	border-style: solid;
	/*
	Hack for Safari browser as button element created from uppy cannot be a flex container, so it's not aligning items vertically the proper way.
	*/
@media not all and (min-resolution:.001dpcm) { 
	@media {
		 div{
		  padding: 0;
		  }
	  }
}
 svg {
	display: none;
}

.uppy-DragDrop-label {
	width: 318px;
	height: 17px;
	font-family: 'Poppins Medium';
	font-size: 12px;
	line-height: 1.64;
	text-align: center;
}

.uppy-Dashboard-browse {
	color:red;
}

.uppy-Dashboard-AddFiles-info {
	display:none;
}

// .uppy-Dashboard-inner {
// 	width: 358px;
// 	display:flex;
// 	justify-content:center;
// 	align-items:center;
// }
`;

const StyledMarkdown = styled(Markdown)`
  p,
  li {
    font-size: 14px;
  }

  p {
    max-width: 100%;
  }
`;

const StyledStatusBar = styled(StatusBar)`
  background: red;
`;

const StyledSelectLabel = styled(Box)`
  font-size: 14px;
  line-height: 1.5;
  color: #541388;
`;

const StyledDropdDown = styled(Box)`
  button {
    border: none;
    border-bottom: 1px solid #f0f0f3;
    border-radius: 0;
  }

  svg {
    width: 32px;
    height: 32px;
    stroke: #000000;
    fill: #000000;
  }
`;

const StyledDashboard = styled(Dashboard)`
  display: flex !important;
  justify-content: center !important;
  align-items: center !important;
`;

const Content = `
<span style="color:#541388"><strong>CSV format:</strong></span>

Need help? Check out our in-depth guide to importing CSVs 

The first row in your file should contain the column headers:  	
- name  
- email  
- phone_number

\`*\` only name, email, phone_number are required.  
\`**\` duplicate email, phone_numbers will be removed  
\`***\` you may pass any other info you would like
`;

const ImportSubscribers = () => {
  const { createNotification } = useContext(NotificationsContext);
  const [selected, setSelected] = useState([]);
  const [segments, callApi] = useApi(
    {
      url: `/api/segments?per_page=40`,
    },
    {
      collection: [],
      links: {
        next: null,
      },
    }
  );

  // const reducer = (state, action) => {
  // 	switch (action.type) {
  // 		case 'append':
  // 			return [ ...state, ...action.payload ];
  // 		default:
  // 			throw new Error('invalid action type.');
  // 	}
  // };

  const [options, setOptions] = useState([]);

  useEffect(() => {
    if (segments.isError || segments.isLoading) {
      return;
    }

    if (segments && segments.data) {
      setOptions(segments.data.collection);
    }
  }, [segments]);

  const onMore = () => {
    if (segments.isError || segments.isLoading) {
      return;
    }

    let url = "";
    if (segments && segments.data && segments.data.links) {
      url = segments.data.links.next;
    }

    if (!url) {
      return;
    }

    callApi({
      url: url,
    });
  };

  const onChange = ({ value: nextSelected }) => {
    setSelected(nextSelected);
  };

  const uppy = useMemo(() => {
    return new Uppy({
      restrictions: {
        maxNumberOfFiles: 1,
        allowedFileTypes: ["text/csv"],
      },
      // onBeforeFileAdded: function() {
      // 	setimportDisabled(false);
      // }
    });
    // .on('file-added', (file) => {
    // 	console.log('Added file', file);
    // 	setimportDisabled(false);
    // });
  });

  // const uppy = Uppy({
  // 	restrictions: {
  // 		maxNumberOfFiles: 1,
  // 		allowedFileTypes: [ 'text/csv' ]
  // 	}
  // onBeforeFileAdded: function(currentFile) {
  // 	setimportDisabled(false);
  // 	if (isOpen) closeNotification();

  // 	const fname = currentFile.name.toLowerCase();
  // 	if (!fname.endsWith('.csv')) {
  // 		uppy.info(`Wrong file type`, 'error', 500);
  // 		return false;
  // 	}
  // 	console.log(currentFile);
  // 	return currentFile;
  // }
  // onBeforeUpload: (files) => {
  // 	// We’ll be careful to return a new object, not mutating the original `files`
  // 	console.log('onBeforeUpload', files);
  // 	const updatedFiles = {};
  // 	Object.keys(files).forEach((fileID) => {
  // 		updatedFiles[fileID] = {
  // 			...files[fileID],
  // 			name: 'myCustomPrefix' + '__' + files[fileID].name
  // 		};
  // 	});
  // 	return updatedFiles;
  // }
  // });

  // uppy.on('file-added', (file) => {
  // 	console.log(file);
  // 	setimportDisabled(false);
  // });

  // uppy.on('file-added', (file) => {
  // 	console.log('Added file', file);
  // 	setimportDisabled(false);
  // });
  // uppy - DragDrop - container;

  uppy.use(AwsS3, {
    async getUploadParameters(file) {
      try {
        const res = await axios.post(
          "/api/s3/sign",
          qs.stringify({
            filename: file.name,
            content_type: file.type,
            action: "import",
          })
        );

        console.log(res);

        return res.data;
      } catch (error) {
        console.log(error);
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
        "/api/subscribers/import",
        qs.stringify(
          {
            filename: file.name,
            segments: selected.map((s) => s.id),
          },
          { arrayFormat: "brackets" }
        )
      );

      createNotification(res.data.message, "status-ok");
    } catch (error) {
      let msg = "Unable to import subscribers. Please try again.";
      if (error.response) {
        msg = error.response.data.message;
      }

      createNotification(msg, "status-error");
    }

    uppy.reset();
  });

  uppy.on("upload-error", async (file, error) => {
    let msg = "Unable to import subscribers. Please try again.";
    if (error.response) {
      msg = error.response.data.message;
    }

    createNotification(msg, "status-error");
  });

  useEffect(() => {
    return () => {
      uppy.close();
    };
  }, [uppy]);

  return (
    <Box direction="column">
      <Box direction="row" animation="fadeIn">
        <Box pad="20px" style={{ minHeight: "auto" }} background=" #fadcff">
          <StyledMarkdown>{Content}</StyledMarkdown>
          <Box
            direction="row"
            width="243px"
            round
            height="39px"
            background="#541388"
            pad={{ top: "10px", bottom: "10px", left: "20px", right: "20px" }}
            style={{ fontSize: "14px", cursor: "pointer" }}
            align="center"
            justify="center"
          >
            <FontAwesomeIcon icon={faDownload} />
            <Box pad={{ left: "7px" }}>{"Download sample .csv file"}</Box>
          </Box>
        </Box>
        <Box
          direction="column"
          pad={{ vertical: "0", left: "14px", right: "20px" }}
        >
          <StyledDropdDown>
            <StyledSelectLabel margin={{ top: "20px", bottom: "10px" }}>
              {" "}
              Add to Group ( Optional )
            </StyledSelectLabel>
            <Select
              multiple
              closeOnChange={false}
              placeholder="Select Group"
              value={selected}
              labelKey="name"
              valueKey="id"
              options={options}
              dropHeight="medium"
              onMore={onMore}
              onChange={onChange}
            />
          </StyledDropdDown>
          <StyledDragDrop
            height="107px"
            width="358px"
            justify="center"
            align="center"
            margin={{ top: "65px", right: "31px", bottom: "0", left: "30px" }}
            border={{ color: "black", size: "medium" }}
          >
            {" "}
            {/* <DragDrop
							width="418px"
							height="107px"
							uppy={uppy}
							locale={{
								strings: {
									dropHereOr: "Drag 'n' drop some files here, or %{browse} to select files",
									browse: 'click'
								}
							}}
						/> */}
            <StyledDashboard
              width="358px"
              height="107px"
              uppy={uppy}
              note={null}
              hideProgressAfterFinish
              showSelectedFiles
              inline
              hideCancelButton
              locale={{
                strings: {
                  dropPasteFiles:
                    "Drag 'n' drop some files here, or %{browse} to select files",
                  browse: "click",
                  dropHint: "",
                },
              }}
              // target={DragDrop}
              // replaceTargetContent
              //  uppy={uppy}
            />
          </StyledDragDrop>{" "}
          <StyledStatusBar
            hideAfterFinish={false}
            showProgressDetails={false}
            uppy={uppy}
            hideUploadButton
          />
        </Box>
      </Box>

      <StyledImportButton
        background="#f5f5fa"
        justify="center"
        align="center"
        pad={{ vertical: "15px" }}
      >
        <Box
          direction="row"
          onClick={() => uppy.upload()}
          justify="center"
          align="center"
          width="144px"
          height="39px"
        >
          <FontAwesomeIcon icon={faCloudUploadAlt} />
          <Text margin={{ left: "10px" }}>Import</Text>
        </Box>
      </StyledImportButton>
    </Box>
  );
};

export default ImportSubscribers;
