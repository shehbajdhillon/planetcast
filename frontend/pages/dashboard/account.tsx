import Navbar from "@/components/dashboard/navbar";
import { Box, HStack, Skeleton, Tab, TabList, TabPanel, TabPanels, Tabs, Text, useColorModeValue } from "@chakra-ui/react";
import Head from "next/head";

const AccountSettings: React.FC = () => {

  const textColor = useColorModeValue("black", "white");
  const bgColor = useColorModeValue("white", "black");

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
      <Box position={"fixed"} top={0} left={0} w="full" p="10px" backgroundColor={useColorModeValue("white", "black")} zIndex={1000}>
        <Navbar />
      </Box>

      <Box pt="80px">

        <Tabs
          variant="enclosed"
          colorScheme="gray"
          isLazy
        >
          <Box
            position={'fixed'}
            w={'100%'}
            zIndex={'999'}
            display={"flex"}
            alignItems={"center"}
            justifyContent={"center"}
          >
            <TabList pl={'25px'} w="full" maxW={"1920px"} backgroundColor={bgColor}>

              <Tab hidden={false}>
                <Text textColor={textColor}>
                  Account Settings
                </Text>
              </Tab>

              <HStack hidden={true} py="10px" spacing="15px">
                <Skeleton h="42px" w="100px" rounded={"md"} />
              </HStack>

            </TabList>
          </Box>

          <TabPanels overflow={'auto'} pt={10}>
            <TabPanel>
              <div />
            </TabPanel>
          </TabPanels>
        </Tabs>
      </Box>

    </Box>
  );
};

export default AccountSettings;
