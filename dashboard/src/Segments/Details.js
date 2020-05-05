import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Grid, Box, Heading, Text, Meter, Button, Paragraph } from "grommet";
import {
  Edit,
  Trash,
  FormClose,
  FormPreviousLink,
  FormNextLink,
} from "grommet-icons";

import { useApi } from "../hooks";
import {
  LoadingOverlay,
  Modal,
  SecondaryButton,
  PlaceholderTable,
} from "../ui";
import EditSegment from "./Edit";
import DeleteSegment from "./Delete";
import RemoveSubscriber from "./RemoveSubscriber";
import history from "../history";
import { Table, Header } from "../Subscribers";

const SubscribersInfoBox = React.memo(({ totalInSegment, total }) => (
  <>
    <Box
      alignSelf="start"
      round={{ corner: "top", size: "small" }}
      background="white"
      pad={{ vertical: "small", right: "large" }}
    >
      <Text margin={{ left: "small" }} size="large">
        <strong>Subscribers</strong>
      </Text>
      <Text size="large" margin={{ top: "small", left: "small" }}>
        <strong>{totalInSegment}</strong>
      </Text>
      <Meter
        round
        max={total}
        margin={{ top: "small", left: "small" }}
        values={[
          {
            color: "brand",
            value: totalInSegment,
            label: "subscribers meter",
          },
        ]}
        aria-label="subscribers meter"
      />
    </Box>
    <Box margin="small">
      <Text>out of {total} total</Text>
    </Box>
  </>
));

SubscribersInfoBox.displayName = "SubscribersInfoBox";
SubscribersInfoBox.propTypes = {
  totalInSegment: PropTypes.number,
  total: PropTypes.number,
};

const Details = ({ match }) => {
  const [segment, setSegment] = useState();
  const [showEdit, setShowEdit] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showDeleteSub, setShowDeleteSub] = useState({ show: false, id: "" });

  const [state, callSegApi] = useApi(
    {
      url: `/api/segments/${match.params.id}`,
    },
    {
      subscribers_in_segment: 0,
      total_subscribers: 0,
    }
  );

  const [subscribers, callApi] = useApi(
    {
      url: `/api/segments/${match.params.id}/subscribers`,
    },
    {
      total: 0,
      collection: [],
    }
  );

  useEffect(() => {
    if (state.isLoading || state.isError) {
      return;
    }

    setSegment(state.data);
  }, [state]);

  if (state.isLoading) {
    return <LoadingOverlay />;
  }

  if (state.isError) {
    return (
      <Box margin="15%" alignSelf="center">
        <Heading>Segment not found.</Heading>
      </Box>
    );
  }

  let table = null;
  if (subscribers.isLoading) {
    table = <PlaceholderTable header={Header} numCols={3} numRows={3} />;
  } else if (subscribers.data.collection.length > 0) {
    table = (
      <Table
        actions={(subscriber) => (
          <Button
            plain
            icon={<FormClose />}
            onClick={() => setShowDeleteSub({ show: true, id: subscriber.id })}
          />
        )}
        list={subscribers.data.collection}
      />
    );
  }

  return (
    <Grid
      rows={["xsmall", "1fr", "1fr"]}
      columns={["260px", "700px", "xsmall"]}
      margin="medium"
      gap="small"
      areas={[
        ["title", "title", "title"],
        ["info", "main", "main"],
        ["info", "main", "main"],
      ]}
    >
      {segment && (
        <>
          {showEdit && (
            <Modal
              title={`Edit segment`}
              hideModal={() => setShowEdit(false)}
              form={
                <EditSegment
                  segment={segment}
                  setSegment={setSegment}
                  hideModal={() => setShowEdit(false)}
                />
              }
            />
          )}
          {showDelete && (
            <Modal
              title={`Delete segment ${segment.name} ?`}
              hideModal={() => setShowDelete(false)}
              form={
                <DeleteSegment
                  id={segment.id}
                  onSuccess={() => history.replace("/dashboard/segments")}
                  onCancel={() => setShowDelete(false)}
                />
              }
            />
          )}
          {showDeleteSub.show && (
            <Modal
              title={`Remove subscriber from segment ?`}
              hideModal={() => setShowDeleteSub({ show: false, id: "" })}
              form={
                <RemoveSubscriber
                  id={segment.id}
                  subId={showDeleteSub.id}
                  onSuccess={async () => {
                    await callApi({
                      url: `/api/segments/${match.params.id}/subscribers`,
                    });
                    await callSegApi({
                      url: `/api/segments/${match.params.id}`,
                    });
                    setShowDeleteSub({ show: false, id: "" });
                  }}
                  onCancel={() => setShowDeleteSub({ show: false, id: "" })}
                />
              }
            />
          )}
          <Box gridArea="title" direction="row">
            <Heading level="2" alignSelf="center">
              {segment.name}
            </Heading>
            <Box direction="row" margin={{ left: "auto" }}>
              <SecondaryButton
                margin={{ right: "small" }}
                a11yTitle="edit segment name"
                alignSelf="center"
                icon={<Edit a11yTitle="edit segment name" color="dark-1" />}
                label="Edit"
                onClick={() => setShowEdit(true)}
              />
              <SecondaryButton
                a11yTitle="delete segment"
                alignSelf="center"
                icon={<Trash a11yTitle="delete segment" color="dark-1" />}
                label="Delete"
                onClick={() => setShowDelete(true)}
              />
            </Box>
          </Box>
          <Box gridArea="info" direction="column">
            <SubscribersInfoBox
              totalInSegment={segment.subscribers_in_segment}
              total={segment.total_subscribers}
            />
          </Box>
          <Box gridArea="main" margin={{ left: "small" }}>
            <Box>
              {table}
              {!subscribers.isLoading &&
              subscribers.data.collection.length === 0 ? (
                <Box align="center" margin={{ top: "small" }}>
                  <Box align="start">
                    <Heading level="2" margin="none">
                      Segment is empty.
                    </Heading>
                  </Box>
                </Box>
              ) : null}
            </Box>
            {!subscribers.isLoading &&
            !subscribers.isError &&
            subscribers.data.collection.length > 0 ? (
              <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
                <Box margin={{ right: "small" }}>
                  <Button
                    icon={<FormPreviousLink />}
                    label="Previous"
                    disabled={subscribers.data.links.previous === null}
                    onClick={() => {
                      callApi({
                        url: subscribers.data.links.previous,
                      });
                    }}
                  />
                </Box>
                <Box>
                  <Button
                    icon={<FormNextLink />}
                    reverse
                    label="Next"
                    disabled={subscribers.data.links.next === null}
                    onClick={() => {
                      callApi({
                        url: subscribers.data.links.next,
                      });
                    }}
                  />
                </Box>
              </Box>
            ) : null}
          </Box>
        </>
      )}
    </Grid>
  );
};

Details.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      id: PropTypes.string,
    }),
  }),
};

export default Details;
