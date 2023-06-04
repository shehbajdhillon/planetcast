import {
  Box,
  Button,
  Grid,
  GridItem,
  useColorModeValue,
  Center,
  useBreakpointValue,
  Text,
} from "@chakra-ui/react";
import { NextPage } from "next";

import { UserResource } from '@clerk/types';

import Image from "next/image";
import { LayoutDashboard } from "lucide-react";
import { Dispatch, SetStateAction } from "react";
import useWindowDimensions from "@/hooks/useWindowDimensions";
import DashboardTab from "@/components/dashboard/dashboard_tab";
import Head from "next/head";
import { MenuBar } from "@/components/dashboard/navbar";


interface SidebarProps {
  user: UserResource | null | undefined;
  signOut: () => void;
  tabIdx: number;
  setTabIdx: Dispatch<SetStateAction<number>>;
};

const Sidebar: React.FC<SidebarProps> = (props) => {

  const { tabIdx, setTabIdx } = props;

  const imageSize = useBreakpointValue({ base: 40, lg: 60 });

  return (
    <Box
      display={"flex"}
      flexDir={{ base: "row", md: "column" }}
      height={"full"}
      w={"full"}
      alignItems={{ base: "center", lg: "inherit"}}
      justifyContent={{ base: "space-between", md: "inherit" }}
    >
      <Center mb={{ base: "0px", md: "20px" }}>
        <Image
          src={useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg')}
          width={imageSize}
          height={100}
          alt='planet cast logo'
        />
      </Center>
      <Button
        my={{ base: "0px", md: "2px" }}
        textAlign={"left"}
        display="flex"
        size={{ base: "md", lg: "lg"}}
        justifyContent="flex-start"
        borderWidth={tabIdx === 0 ? "1px" : "0px"}
        variant={tabIdx === 0 ? "solid" : "ghost"}
        onClick={() => setTabIdx(0)}
      >
        <LayoutDashboard />
        <Text pl={"10px"} display={{ base:"none", lg: "flex" }}>
          Dashboard
        </Text>
      </Button>
      <Box marginTop={{ base: "0px", md: "auto"}} marginBottom={"10px"}>
        <Center>
          <MenuBar />
        </Center>
      </Box>
    </Box>
  );
};

const Dashboard: NextPage = () => {

  const { height } = useWindowDimensions();

  return (
    <Box>
      <Head>
        <title>Dashboard | PlanetCast</title>
        <meta
          name="description"
          content="Cast Content in any Language, Across the Planet"
        />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Box
        display={"flex"}
        justifyContent={"center"}
      >
        <Grid
          h={height}
          templateAreas={{
            base: `
              "main"
            `,
          }}
          gridTemplateColumns={{
            base: "1fr",
          }}
          maxW={"1920px"}
          w="full"
        >
          <GridItem
            area={"main"}
            display={"flex"}
            flexDir={"column"}
          >
            <DashboardTab />
          </GridItem>
        </Grid>
      </Box>
    </Box>
  );
};

export default Dashboard;