import {
  Box,
  Button,
  Badge,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent,
  DrawerHeader,
  DrawerOverlay,
  Grid,
  GridItem,
  HStack,
  Heading,
  Spacer,
  Stack,
  Text,
  VStack,
  useDisclosure,
  useToast
} from "@chakra-ui/react";
import { Menu } from "lucide-react";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import StripeCheckoutForm from "../stripe_checkout";

import { useSearchParams } from 'next/navigation';
import { Team } from "@/types";

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
        General
      </Button>
      <Button
        variant={"ghost"}
        w="full"
        onClick={() => switchTab(1)}
        borderWidth={tabIdx === 1 ? '1px' : ''}
        textAlign={"left"}
        justifyContent={"left"}
      >
        Subscription
      </Button>
    </VStack>
  );
};

interface SettingsTabProps {
  teamSlug: string;
  params: string[];
  refetch: () => void;
  loading: boolean;
  team: Team
};

const SettingsTab: React.FC<SettingsTabProps> = ({ teamSlug, params, loading, team }) => {

  const [tabIdx, setTabIdx] = useState(getInitialTeamSettingTabIndex(params).defaultIndex);

  const router = useRouter();

  const { isOpen, onOpen, onClose } = useDisclosure();

  useEffect(() => {
    const team = router.query.team as string[];
    const tabProps = getInitialTeamSettingTabIndex(team);
    setTabIdx(tabProps.defaultIndex);
  }, [router.query]);

  const switchTab = (tabIdx: number) => {
    setTabIdx(tabIdx);
    switch(tabIdx) {
      case 0:
        router.push(`/dashboard/${teamSlug}/settings`, undefined, { shallow: true });
        break;
      case 1:
        router.push(`/dashboard/${teamSlug}/settings/subscription`, undefined, { shallow: true });
        break;
    }
  };

  return (
    <Box
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      w={"full"}
      hidden={loading}
    >
      <Box w="full" maxW={"1200px"}>
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
          <GridItem area={"main"} maxW={"912px"}>

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

            {tabIdx === 0 && <GeneralSettingsTab drawerOpen={onOpen} />}
            {tabIdx === 1 && <SubscriptionSettingsTab drawerOpen={onOpen} teamSlug={teamSlug} team={team} />}
          </GridItem>
        </Grid>
      </Box>
    </Box>
  );
};

interface SubscriptionSettingsTabProps {
  drawerOpen: () => void;
  teamSlug: string;
  team: Team
};

const SubscriptionSettingsTab: React.FC<SubscriptionSettingsTabProps> = (props) => {
  const { drawerOpen, teamSlug, team } = props;

  const searchParams = useSearchParams();
  const router = useRouter();
  const toast = useToast();

  const currentSubscription = team?.subscriptionPlans[0];

  useEffect(() => {
    const action = searchParams.get("action");

    if (!action) return;

    if (action === "cancel") {
      toast({
        title: 'Subscription Update Cancelled',
        description: "No change has been made to your current subscription.",
        status: 'info',
        duration: 6000,
        isClosable: true,
        position: "top",
        containerStyle: {
          paddingTop: "30px"
        },
      });
    }

    if (action === "success") {
      toast({
        title: 'Subscription Updated',
        description: "Your subscription has been successfully updated.",
        status: 'success',
        duration: 6000,
        isClosable: true,
        position: "top",
        containerStyle: {
          paddingTop: "30px"
        },
      });
    }

    router.replace(`/dashboard/${teamSlug}/settings/subscription`, undefined, { shallow: true });

  }, []);

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
          <Text>Subscription</Text>
        </HStack>
      </Button>
      <Heading>Subscription</Heading>
      <Stack
        direction={"column"}
        w="full"
        borderWidth={"1px"}
        padding={"25px"}
        rounded={"lg"}
      >
        <VStack h="full" w="full" spacing="15px">
          {currentSubscription?.subscriptionActive ?
            <HStack w="full">
              <Badge colorScheme="green" rounded={"md"} p={1}>
                Active
              </Badge>
            <Text fontWeight={"semibold"}>Starter Monthly</Text>
            </HStack>
          :
            <HStack w="full">
              <Badge colorScheme="gray" rounded={"md"} p={1}>
                Inactive
              </Badge>
            <Text fontWeight={"semibold"}>No Active Subscription</Text>
            </HStack>
          }


          {currentSubscription?.subscriptionActive &&
            <>
              <Stack w="full" spacing={"2px"}>
                <Text fontWeight={"semibold"}>Current Billing Cycle:</Text>
                <Text fontWeight={"semibold"}>Nov 15 2023 to Dec 25 2023</Text>
              </Stack>

              <Stack w="full" spacing={"2px"}>
                <Text fontWeight={"semibold"}>Next Billing Cycle:</Text>
                <Text fontWeight={"semibold"}>Dec 25 2023 to Jan 25 2023</Text>
              </Stack>
            </>
          }

          <Stack w="full" spacing={"2px"}>
            <Text fontWeight={"semibold"}>Remaining Minutes:</Text>
            <Text fontWeight={"semibold"}>{currentSubscription?.remainingCredits}</Text>
          </Stack>

          {currentSubscription?.subscriptionActive &&
            <Stack w="full" spacing={"2px"}>
              <Text fontWeight={"semibold"}>Subscription automatically renews on Dec 25 2023</Text>
              <Text>Card ending with 0045 will be charged $60</Text>
              <Text>Click on Manage Billing to update payment information</Text>
            </Stack>
          }

          <Box w="full" display="flex">
            <StripeCheckoutForm teamSlug={teamSlug} subscriptionActive={currentSubscription?.subscriptionActive}/>
          </Box>

        </VStack>
      </Stack>
    </VStack>
  );
};

interface GeneralSettingsTabProps {
  drawerOpen: () => void;
};

const GeneralSettingsTab: React.FC<GeneralSettingsTabProps> = (props) => {

  const { drawerOpen } = props;

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
          <Text>General</Text>
        </HStack>
      </Button>
      <Heading>Danger Zone</Heading>
      <Stack
        direction={"column"}
        w="full"
        borderColor={"red.200"}
        borderWidth={"1px"}
        padding={"25px"}
        rounded={"lg"}
      >
        <HStack>
          <Box>
            <Text>Delete this Team</Text>
            <Text>
              This team and all the projects created will be deleted. This action is not reversible.
            </Text>
          </Box>
          <Spacer />
          <Box>
            <Button isDisabled={true} colorScheme="red">
              Delete Team
            </Button>
          </Box>
        </HStack>
      </Stack>
    </VStack>

  );
};

export const getInitialTeamSettingTabIndex = (tabs: string[]) => {

  let defaultIndex = 0;
  let notFound = false;

  if (!tabs || tabs.length < 2) {
    return {
      defaultIndex,
      notFound,
    } as const;
  }

  const [_, tabName] = tabs;

  switch(tabName) {
    case undefined:
      defaultIndex = 0;
      break;
    case 'subscription':
      defaultIndex = 1;
      break;
    default:
      notFound = true
      break;
  }

  return {
    defaultIndex,
    notFound,
  } as const;

};

export default SettingsTab;
