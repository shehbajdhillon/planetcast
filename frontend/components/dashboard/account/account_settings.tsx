import { Team, TeamInvite } from "@/types";
import { HStack, Button, Text, Heading, Stack, VStack, Box, Spacer, Grid, GridItem, Drawer, DrawerOverlay, DrawerContent, DrawerCloseButton, DrawerHeader, DrawerBody, useDisclosure } from "@chakra-ui/react";
import { Menu } from "lucide-react";
import { useState } from "react";

interface TabButtonsProps {
  tabIdx: number;
  switchTab: (tabIdx: number) => void;
};

const TabButtons: React.FC<TabButtonsProps> = ({ tabIdx, switchTab }) => {
  return (
    <VStack w="full" alignItems={"flex-start"} px="10px" spacing={"10px"}>
      <Button
        w="full"
        variant={"ghost"}
        onClick={() => switchTab(0)}
        borderWidth={tabIdx === 0 ? '1px' : ''}
        justifyContent={"left"}
      >
        Teams
      </Button>
    </VStack>
  );
};

interface TeamInvitesTabProps {
  drawerOpen: () => void;
  teams: any[];
  invites: any[];
}

const TeamInvitesTab: React.FC<TeamInvitesTabProps> = (props) => {
  const { teams, invites, drawerOpen } = props;

  return (
    <VStack alignItems={{ lg: "flex-start" }}>
      <Button
        w="full"
        variant={"ghost"}
        borderWidth={'1px'}
        display={{ lg: "none" }}
        textAlign={"left"}
        onClick={drawerOpen}
      >
        <HStack>
          <Menu size={"20px"} />
          <Text>Teams</Text>
        </HStack>
      </Button>
      <Heading>Teams</Heading>
      <Stack
        direction={"column"}
        w="full"
        borderWidth={"1px"}
        padding={"25px"}
        rounded={"lg"}
      >
        <HStack>
          <Box>
            <Text>Current Teams</Text>
            {teams?.map((team, idx) => (
              <Text>{team.teamName}</Text>
            ))}
          </Box>
          <Spacer />
          <Box>
          </Box>
        </HStack>
      </Stack>
      <Stack
        direction={"column"}
        w="full"
        borderWidth={"1px"}
        padding={"25px"}
        rounded={"lg"}
      >
        <HStack>
          <Box>
            <Text>Current Invites</Text>
            {invites?.map((invite, idx) => (
              <Text>{invite.teamName}</Text>
            ))}
          </Box>
          <Spacer />
          <Box>

          </Box>
        </HStack>
      </Stack>
    </VStack>
  );
};

interface AccountSettingsTab {
  loading: boolean;
  teams: Team[];
  invites: TeamInvite[];
};

const AccountSettingsTab: React.FC<AccountSettingsTab> = (props) => {

  const { teams, invites, loading } = props;
  const { isOpen, onOpen, onClose } = useDisclosure();

  const [tabIdx, setTabIdx] = useState(0);
  const switchTab = (tabIdx: number) => {
    setTabIdx(tabIdx);
  };

  return (
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
      hidden={loading}
    >
      <Box w="full" maxW={"1450px"}>
        <Grid
          templateAreas={{
            base: `
              "main"
            `,
            lg: `"sidebar main"`
          }}
          gridTemplateColumns={{ base: "1fr", lg: "1fr 4fr"}}
          w="full"
          h="full"
          gap="10px"
        >
          <GridItem area={"sidebar"} display={{ base: "none", lg: "block" }}>
            <TabButtons tabIdx={tabIdx} switchTab={switchTab} />
          </GridItem>
          <GridItem area={"main"} w="full">

            <Drawer isOpen={isOpen} onClose={onClose} placement="left" size={"md"}>
              <DrawerOverlay />
              <DrawerContent>
                <DrawerCloseButton />
                <DrawerHeader>Team Settings</DrawerHeader>
                <DrawerBody>
                  <TabButtons
                    switchTab={(idx: number) => {
                      switchTab(idx);
                      onClose();
                    }}
                    tabIdx={tabIdx}
                  />
                </DrawerBody>
              </DrawerContent>
            </Drawer>

            {tabIdx === 0 && <TeamInvitesTab drawerOpen={onOpen} teams={teams} invites={invites} />}
          </GridItem>
        </Grid>
      </Box>
    </Box>
  );
};

export default AccountSettingsTab;
