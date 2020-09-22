import React, { useState, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import { Box, Heading, ResponsiveContext, Grid, Text, Meter } from "grommet";

import { LoadingOverlay } from "../ui";
import { useApi } from "../hooks";

const DetailsGrid = ({ children }) => {
  const size = useContext(ResponsiveContext);

  let cols = ["small", "small", "large", "xsmall"];
  let areas = [
    [".", "title", "title", "title"],
    [".", "info", "main", "main"],
    [".", "info", "main", "main"],
  ];

  if (size === "medium") {
    cols = ["264px", "600px", "xsmall"];
    areas = [
      ["title", "title", "title"],
      ["info", "main", "main"],
      ["info", "main", "main"],
    ];
  }

  return (
    <Grid
      rows={["xsmall", "1fr", "1fr"]}
      columns={cols}
      margin="medium"
      gap="small"
      areas={areas}
    >
      {children}
    </Grid>
  );
};

DetailsGrid.displayName = "DetailsGrid";
DetailsGrid.propTypes = {
  children: PropTypes.element,
};

const Stat = React.memo(({ label, value, total, footer }) => (
  <>
    <Box direction="row">
      <Text margin={{ right: "large" }} size="medium">
        <b>{label}</b>
      </Text>
      <Text margin={{ left: "auto" }} size="medium">
        {Math.round((value / total) * 100).toFixed(2)}%
      </Text>
    </Box>
    <Box align="end" margin={{ left: "small", top: "xsmall" }}>
      <Text size="16px">{footer}</Text>
    </Box>
  </>
));

Stat.displayName = "Stat";

const PercentageStats = React.memo(({ stats }) => {
  const s = [
    {
      label: "Delivered",
      value: stats.delivered,
      total: stats.total_sent,
      footer: `${stats.delivered} out of ${stats.total_sent} total`,
    },
    {
      label: "Opens",
      value: stats.opens.total,
      total: stats.total_sent,
      footer: `${stats.opens.unique} unique ${stats.opens.total} total`,
    },
    {
      label: "Clicks",
      value: stats.clicks.total,
      total: stats.total_sent,
      footer: `${stats.clicks.unique} unique ${stats.clicks.total} total`,
    },
    {
      label: "Bounces",
      value: stats.bounces,
      total: stats.total_sent,
      footer: `${stats.bounces} bounced`,
    },
    {
      label: "Complaints",
      value: stats.complaints,
      total: stats.total_sent,
      footer: `${stats.complaints} complained`,
    },
  ];
  return (
    <>
      <Box
        alignSelf="start"
        round={{ corner: "top", size: "small" }}
        background="white"
        pad={{ vertical: "small", right: "medium" }}
      >
        {s.map((stat) => (
          <Box
            key={stat.label}
            margin={{ bottom: "medium" }}
            pad={{ horizontal: "small" }}
          >
            <Stat {...stat} />
          </Box>
        ))}
      </Box>
    </>
  );
});

PercentageStats.displayName = "PercentageStats";
PercentageStats.propTypes = {
  stats: PropTypes.shape({
    total_sent: PropTypes.number,
    delivered: PropTypes.number,
    opens: PropTypes.shape({
      unique: PropTypes.number,
      total: PropTypes.number,
    }),
    clicks: PropTypes.shape({
      unique: PropTypes.number,
      total: PropTypes.number,
    }),
    bounces: PropTypes.number,
    complaints: PropTypes.number,
  }),
};

const Details = ({ match }) => {
  const [campaign, setCampaign] = useState();

  const [state] = useApi({
    url: `/api/campaigns/${match.params.id}`,
  });

  const [stats] = useApi({
    url: `/api/campaigns/${match.params.id}/stats`,
  });

  useEffect(() => {
    if (state.isLoading || state.isError) {
      return;
    }

    setCampaign(state.data);
  }, [state]);

  if (state.isLoading) {
    return <LoadingOverlay />;
  }

  if (state.isError) {
    return (
      <Box margin="15%" alignSelf="center">
        <Heading>Campaign not found.</Heading>
      </Box>
    );
  }

  return (
    <DetailsGrid>
      {campaign && (
        <>
          <Box gridArea="title" direction="row">
            <Heading level="2" alignSelf="center">
              {campaign.name} stats report
            </Heading>
          </Box>
          <Box gridArea="info" direction="column">
            {stats && stats.data && <PercentageStats stats={stats.data} />}
          </Box>
        </>
      )}
    </DetailsGrid>
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
