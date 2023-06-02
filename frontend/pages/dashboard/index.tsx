import {
  Box,
  Button,
  Grid,
  GridItem,
  useColorModeValue,
  Menu,
  MenuButton,
  MenuList,
  Avatar,
  Center,
  MenuDivider,
  MenuItem,
  useBreakpointValue,
  Text,
  HStack,
  useColorMode,
} from "@chakra-ui/react";
import { NextPage } from "next";

import { useClerk, useUser } from "@clerk/nextjs";
import { UserResource } from '@clerk/types';

import Image from "next/image";
import { LayoutDashboard, Orbit, } from "lucide-react";
import { Dispatch, SetStateAction, useState } from "react";
import useWindowDimensions from "@/hooks/useWindowDimensions";
import DashboardTab from "@/components/dashboard/dashboard_tab";
import AccountsTab from "@/components/dashboard/planetcast_tab";
import Head from "next/head";

interface MenuBarProps {
  fullName: string;
  emailAddress: string;
  logout: () => void;
};

const MenuBar: React.FC<MenuBarProps> = ({ emailAddress, fullName, logout }) => {

  const { toggleColorMode } = useColorMode();

  return (
    <Menu>
      <MenuButton
        rounded={'full'}
        cursor={'pointer'}
        borderWidth={"1px"}
        borderColor="blackAlpha.400"
        _hover={{
          backgroundColor: useColorModeValue("blackAlpha.200", "whiteAlpha.200")
        }}
      >
        <HStack>
          <Avatar
            size={{ base: "sm", md: "md", lg: 'lg'}}
            src={`https://api.dicebear.com/6.x/notionists/svg?seed=${emailAddress}`}
            borderRadius={"full"}
            backgroundColor={"white"}
          />
        </HStack>
      </MenuButton>
      <MenuList alignItems={'center'} backgroundColor={useColorModeValue("white", "black")}>
        <Center>
          <Avatar
            size={'2xl'}
            src={`https://api.dicebear.com/6.x/notionists/svg?seed=${emailAddress}`}
            backgroundColor={"white"}
          />
        </Center>
        <Center maxW={"220px"}>
          <Text
            display={{ base: "none", lg: "inline-block" }}
            whiteSpace={"nowrap"}
            textOverflow={"ellipsis"}
            overflow={'hidden'}
            noOfLines={1}
          >
            {fullName}
          </Text>
        </Center>
        <MenuDivider />
        <MenuItem
          backgroundColor={useColorModeValue("white", "black")}
          _hover={{
            backgroundColor: useColorModeValue("blackAlpha.200", "whiteAlpha.200")
          }}
        >
          Settings
        </MenuItem>
        <MenuItem
          backgroundColor={useColorModeValue("white", "black")}
          _hover={{
            backgroundColor: useColorModeValue("blackAlpha.200", "whiteAlpha.200")
          }}
          onClick={toggleColorMode}
        >
          Turn on {useColorModeValue("Dark", "Light")} Mode
        </MenuItem>
        <MenuItem
          backgroundColor={useColorModeValue("white", "black")}
          onClick={logout}
          _hover={{
            backgroundColor: useColorModeValue("blackAlpha.200", "whiteAlpha.200")
          }}
        >
          Logout
        </MenuItem>
      </MenuList>
    </Menu>
  )
};

interface SidebarProps {
  user: UserResource | null | undefined;
  signOut: () => void;
  tabIdx: number;
  setTabIdx: Dispatch<SetStateAction<number>>;
};

const Sidebar: React.FC<SidebarProps> = (props) => {

  const { user, signOut, tabIdx, setTabIdx } = props;

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
      <Button
        my={{ base: "0px", md: "2px" }}
        textAlign={"left"}
        size={{ base: "md", lg: "lg"}}
        display="flex"
        justifyContent="flex-start"
        borderWidth={tabIdx === 1 ? "1px" : "0px"}
        variant={tabIdx === 1 ? "solid" : "ghost"}
        onClick={() => setTabIdx(1)}
      >
        <Orbit />
        <Text pl={"10px"} display={{ base:"none", lg: "flex" }}>
          PlanetCast
        </Text>
      </Button>
      <Box marginTop={{ base: "0px", md: "auto"}}>
        <MenuBar
          emailAddress={user?.primaryEmailAddress?.emailAddress || ""}
          fullName={user?.fullName || ""}
          logout={signOut}
        />
      </Box>
    </Box>
  );
};

const Dashboard: NextPage = () => {
  const { user } = useUser();
  const { signOut } = useClerk();
  const [tabIdx, setTabIdx] = useState(0);

  const mobileView = useBreakpointValue({ base: true, md: false });

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
            md: `"sidebar main"`
          }}
          gridTemplateColumns={{
            base: "1fr",
            md: "70px 1fr",
            lg: `220px 1fr`,
          }}
          gridTemplateRows={{
            base: "1fr",
          }}
          maxW={"1920px"}
          w="full"
        >
          <GridItem
            area={"sidebar"}
            borderRightWidth={{ base: "0px", md: "1px" }}
            p="10px"
            display={{ base: "flex" }}
            hidden={mobileView}
          >
            <Sidebar user={user} signOut={signOut} tabIdx={tabIdx} setTabIdx={setTabIdx} />
          </GridItem>

          <GridItem
            area={"main"}
            display={"flex"}
            flexDir={"column"}
          >
            <Box w="full" h="full">
              { tabIdx == 0 && <DashboardTab /> }
              { tabIdx == 1 && <AccountsTab /> }
            </Box>
            <Box
              marginTop={"auto"}
              w="full"
              display={{ md: "none" }}
              borderTopWidth={{ base: "1px", md: "0px" }}
              p="10px"
            >
              <Sidebar user={user} signOut={signOut} tabIdx={tabIdx} setTabIdx={setTabIdx} />
            </Box>
          </GridItem>
        </Grid>
      </Box>
    </Box>
  );
};

export default Dashboard;
