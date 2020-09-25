import React, { useState, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import { Box, Heading, ResponsiveContext, Grid, Text } from "grommet";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
} from "recharts";

import Bounces from "./Stats/Bounces";
import Complaints from "./Stats/Complaints";
import { LoadingOverlay } from "../ui";
import { useApi } from "../hooks";

const DetailsGrid = ({ children }) => {
  const size = useContext(ResponsiveContext);

  let cols = ["small", "small", "large", "xsmall"];
  let areas = [
    [".", "title", "title", "title"],
    [".", "info", "main", "main"],
    [".", "info", "main", "main"],
    [".", "bounces", "bounces", "bounces"],
    [".", "complaints", "complaints", "complaints"],
  ];

  if (size === "medium") {
    cols = ["264px", "600px", "xsmall"];
    areas = [
      ["title", "title", "title"],
      ["info", "main", "main"],
      ["info", "main", "main"],
      ["bounces", "bounces", "bounces"],
      ["complaints", "complaints", "complaints"],
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
  <Box
    margin={{ horizontal: "xsmall" }}
    pad={{ horizontal: "small", vertical: "xsmall" }}
  >
    <Box direction="row">
      <Text margin={{ right: "large" }} size="medium">
        <b>{label}</b>
      </Text>
      <Text margin={{ left: "auto" }} size="medium">
        {total && Math.round((value / total) * 100).toFixed(2)}%
      </Text>
    </Box>
    <Box align="end" margin={{ left: "small", top: "xsmall" }}>
      <Text size="16px">{footer}</Text>
    </Box>
  </Box>
));

Stat.displayName = "Stat";
Stat.propTypes = {
  label: PropTypes.string,
  value: PropTypes.number,
  total: PropTypes.number,
  footer: PropTypes.string,
};

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
      >
        {s.map((stat) => (
          <Box key={stat.label} pad={{ horizontal: "small" }} border="bottom">
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

const StatsChart = React.memo(({ data }) => (
  <BarChart width={700} height={350} data={data}>
    <CartesianGrid strokeDasharray="3 3" />
    <XAxis dataKey="name" />
    <YAxis />
    <Tooltip />
    <Legend />
    <Bar dataKey="unique" stackId="a" fill="#711FFF" />
    <Bar dataKey="total" stackId="a" fill="#00C781" />
  </BarChart>
));

StatsChart.displayName = "StatsChart";
StatsChart.propTypes = {
  data: PropTypes.array,
};

const Details = ({ match }) => {
  const [campaign, setCampaign] = useState();
  const [data, setData] = useState([]);

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

  useEffect(() => {
    if (stats.isLoading || stats.isError || !stats.data) {
      return;
    }

    const data = stats.data;
    setData([
      {
        name: "Total sent",
        total: data.total_sent,
      },
      {
        name: "Delivered",
        total: data.delivered,
      },
      {
        name: "Bounces",
        total: data.bounces,
      },
      {
        name: "Complaints",
        total: data.complaints,
      },
      {
        name: "Clicks",
        total: data.clicks.total,
        unique: data.clicks.unique,
      },
      {
        name: "Opened",
        total: data.opens.total,
        unique: data.opens.unique,
      },
      {
        name: "Unopened",
        total: data.delivered - data.opens.total,
      },
    ]);
  }, [stats]);

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
          <Box gridArea="info">
            {stats && stats.data && <PercentageStats stats={stats.data} />}
          </Box>
          <Box gridArea="main">
            <StatsChart data={data} />
          </Box>
          <Box fill gridArea="bounces">
            <Heading level="3">Bounces</Heading>
            <Bounces campaignId={campaign.id} />
          </Box>
          <Box fill gridArea="complaints" margin={{ bottom: "medium" }}>
            <Heading level="3">Complaints</Heading>
            <Complaints campaignId={campaign.id} />
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
