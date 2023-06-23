import Navbar from "@/components/dashboard/navbar";
import {
  Box,
  useColorModeValue,
} from "@chakra-ui/react";
import { NextPage } from "next";

import Head from "next/head";

const Dashboard: NextPage = () => {
  return (
    <Box>
      <Head>
        <title>Project | lanetCast</title>
        <meta
          name="description"
          content="Cast Content in any Language, Across the Planet"
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box position={"fixed"} top={0} left={0} w="full" p="10px" backgroundColor={useColorModeValue("white", "black")} zIndex={1000}>
        <Navbar />
      </Box>
      <Box
        display={"flex"}
        justifyContent={"center"}
      >
      </Box>
    </Box>
  );
};

export default Dashboard;

