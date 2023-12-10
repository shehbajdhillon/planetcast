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
  useToast,
  useColorModeValue,
  IconButton,
  Avatar,
  Input
} from "@chakra-ui/react";
import {  Menu, UserPlusIcon, UserXIcon } from "lucide-react";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import StripeCheckoutForm from "../stripe_checkout";

import { useSearchParams } from 'next/navigation';
import { Team } from "@/types";
import PricingComponent from "../marketing_page/pricing_component";
import { gql, useMutation } from "@apollo/client";
import { convertUtcToLocal, validateEmail } from "@/utils";
import { loadStripe } from "@stripe/stripe-js";
import useWindowDimensions from "@/hooks/useWindowDimensions";
import SingleActionModal from "../single_action_modal";

const CREATE_STRIPE_CHECKOUT = gql`
  mutation CreateStripeCheckout($teamSlug: String!, $lookUpKey: String!) {
    createCheckoutSession(teamSlug: $teamSlug, lookUpKey: $lookUpKey) {
      sessionId
    }
  }
`;

const SEND_INVITE = gql`
  mutation SendInvite($teamSlug: String!, $inviteeEmail: String!) {
    sendTeamInvite(teamSlug: $teamSlug, inviteeEmail: $inviteeEmail)
  }
`
const DELETE_INVITE = gql`
  mutation DeleteInvite($teamSlug: String!, $inviteSlug: String!) {
    deleteTeamInvite(teamSlug: $teamSlug, inviteSlug: $inviteSlug)
  }
`

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
        Members
      </Button>
      <Button
        variant={"ghost"}
        w="full"
        onClick={() => switchTab(2)}
        borderWidth={tabIdx === 2 ? '1px' : ''}
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

const SettingsTab: React.FC<SettingsTabProps> = ({ refetch, teamSlug, params, loading, team }) => {

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
        router.push(`/dashboard/${teamSlug}/settings/members`, undefined, { shallow: true });
        break;
      case 2:
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

            {tabIdx === 0 && <GeneralSettingsTab drawerOpen={onOpen} />}
            {tabIdx === 1 && <TeamMembersTab drawerOpen={onOpen} team={team} teamSlug={teamSlug} refetch={refetch} />}
            {tabIdx === 2 && <SubscriptionSettingsTab drawerOpen={onOpen} teamSlug={teamSlug} team={team} />}
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
  const [showUpgrade, setShowUpgrade] = useState(false);
  const onUpdateClick = () => {
    setShowUpgrade(true);
    router.push('#upgrades');
  }

  const [annualPricing, setAnnualPricing] = useState(true);
  const bgColor = useColorModeValue("black", "white");
  const textColor = useColorModeValue("white", "black");

  const [createCheckoutSession, { loading }] = useMutation(CREATE_STRIPE_CHECKOUT);

  const handleCheckout = async (lookUpKey: string) => {
    const stripe = await loadStripe(process.env.NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY || "")
    const response = await createCheckoutSession({ variables: { teamSlug, lookUpKey } });
    const sessionId = response.data?.createCheckoutSession?.sessionId;
    const res = await stripe?.redirectToCheckout({ sessionId });
    if (res?.error) {
      console.log('[stripe error]', res.error.message);
    }
  };

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

  const subActive = currentSubscription?.stripeSubscriptionId !== null;
  const subscriptionData = currentSubscription?.subscriptionData

  const subStatus = subscriptionData?.status;

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
          {subActive ?
            <HStack w="full">
              <Badge
                colorScheme={subStatus === "active" ? "green" : "red"}
                rounded={"md"}
                p={1}
              >
                { subStatus }
              </Badge>
            <Text fontWeight={"semibold"}>{ subscriptionData?.planName }</Text>
            </HStack>
          :
            <HStack w="full">
              <Badge colorScheme="gray" rounded={"md"} p={1}>
                Inactive
              </Badge>
            <Text fontWeight={"semibold"}>No Active Subscription</Text>
            </HStack>
          }


          {subActive &&
            <>
              <Stack w="full" spacing={"2px"}>
                <Text fontWeight={"semibold"}>Current Billing Cycle:</Text>
                <Text fontWeight={"semibold"}>
                  {convertUtcToLocal(subscriptionData?.currentPeriodStart || "")} to {' '}
                  {convertUtcToLocal(subscriptionData?.currentPeriodEnd || "")}
                </Text>
              </Stack>
            </>
          }

          <Stack w="full" spacing={"2px"}>
            <Text fontWeight={"semibold"}>Remaining Minutes:</Text>
            <Text fontWeight={"semibold"}>{currentSubscription?.remainingCredits}</Text>
          </Stack>

          {subActive && subStatus === "active" ?
            <Stack w="full" spacing={"2px"}>
              <Text fontWeight={"semibold"}>
                Subscription next renews on {convertUtcToLocal(subscriptionData?.currentPeriodEnd || "")} (renews every {subscriptionData?.interval})
              </Text>
              <Text>Card ending with {subscriptionData?.lastFourCardDigits} will be charged ${subscriptionData?.costInUsd}</Text>
              <Text>Click on Manage Billing to update payment information</Text>
            </Stack>
            :
            <Stack w="full" spacing={"2px"}>
              <Text fontWeight={"semibold"}>
                There is a problem in renewing your subscription.
              </Text>
              <Text fontWeight={"semibold"}>
                Please click on Manage Billing to resolve the issue.
              </Text>
            </Stack>
          }

          <Box w="full" display="flex">
            <StripeCheckoutForm
              teamSlug={teamSlug}
              subscriptionActive={subActive}
              onUpdateClick={onUpdateClick}
            />
          </Box>
        </VStack>
      </Stack>

      <Box id="upgrades" scrollMarginTop={"120px"} />
      { showUpgrade &&
        <VStack spacing={"30px"} pt="20px" pb="50px" w="full">
          <HStack borderWidth={"1px"} p="5px" rounded={"md"}>
            <Button
              onClick={() => setAnnualPricing(false)}
              bgColor={!annualPricing ? bgColor : textColor}
              textColor={!annualPricing ? textColor : bgColor}
              bgGradient={!annualPricing ? 'linear(to-tl, #007CF0, #01DFD8)' : ''}
              _hover={{
                backgroundColor: !annualPricing ? bgColor : textColor,
                textColor: annualPricing ? bgColor : textColor,
              }}
            >
              Monthly
            </Button>
            <Button
              onClick={() => setAnnualPricing(true)}
              bgColor={annualPricing ? bgColor : textColor}
              textColor={annualPricing ? textColor : bgColor}
              bgGradient={annualPricing ? 'linear(to-tl, #007CF0, #01DFD8)' : ''}
              _hover={{
                backgroundColor: annualPricing ? bgColor : textColor,
                textColor: !annualPricing ? bgColor : textColor,
              }}
            >
              Annual
            </Button>
          </HStack>
          <PricingComponent
            p="0px"
            annualPricing={annualPricing}
            spacing={{ base: "16px", md: "4px" }}
            handleCheckout={handleCheckout}
            marketingPage={false}
          />
        </VStack>
      }
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

interface TeamMembersTabProps {
  drawerOpen: () => void;
  team: Team;
  teamSlug: string;
  refetch: () => void;
};

const TeamMembersTab: React.FC<TeamMembersTabProps> = (props) => {

  const { refetch, drawerOpen, teamSlug, team } = props;
  const { width } = useWindowDimensions();
  const { isOpen, onOpen, onClose } = useDisclosure();

  const [sendInvite, { loading, error, data }] = useMutation(SEND_INVITE);
  const [deleteInvite, { loading: deleteLoading, error: deleteError, data: deleteData }] = useMutation(DELETE_INVITE);

  const [inviteeEmail, setInviteeEmail] = useState("");
  const [emailValid, setEmailValid] = useState(false);

  const toast = useToast();

  const sendEmailInvite = async () => {
    const res = await sendInvite({ variables: { teamSlug, inviteeEmail } })
    if (res) refetch();
  };

  const deleteEmailInvite = async (inviteSlug: string) => {
    const res = await deleteInvite({ variables: { teamSlug, inviteSlug } });
    if (res) refetch();
  };

  useEffect(() => setEmailValid(validateEmail(inviteeEmail)), [inviteeEmail]);
  useEffect(() => setInviteeEmail(""), [isOpen]);

  useEffect(() => {
    if (!error) return;
    toast({
      title: 'Could not send invite',
      description: error.message,
      status: 'error',
      duration: 6000,
      isClosable: true,
      position: "top",
      containerStyle: {
        paddingTop: "30px"
      },
    });
  }, [error]);

  useEffect(() => {
    if (!error && !loading && data) {
      toast({
        title: 'Sent invite successfully',
        description: 'Please ask invitee to check their email',
        status: 'success',
        duration: 6000,
        isClosable: true,
        position: "top",
        containerStyle: {
          paddingTop: "30px"
        },
      });
    }
  }, [error, loading, data]);


  return (
    <VStack alignItems={{ lg: "flex-start" }} maxW={"100%"} overflowX={"auto"}>
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
          <Text>Team Members</Text>
        </HStack>
      </Button>
      <HStack w="full">
        <Heading>Manage Members</Heading>
        <Spacer />
        <IconButton
          aria-label="add member"
          variant={"outline"}
          icon={<UserPlusIcon />}
          onClick={onOpen}
        />
        <SingleActionModal
          isOpen={isOpen}
          onClose={onClose}
          heading={"Add Member"}
          loading={loading}
          disabled={!emailValid}
          action={sendEmailInvite}
        >
          <Text>Email Address</Text>
          <Input
            placeholder="john.doe@planetcast.ai"
            onChange={(e) => setInviteeEmail(e.target.value)}
          />
        </SingleActionModal>
      </HStack>

      <Stack
        direction={"column"}
        w="full"
        borderWidth={"1px"}
        padding={"25px"}
        rounded={"lg"}
        spacing={"35px"}
        maxW={(width as number) - 32}
      >

        {team?.members.map((m, idx) => (
          <HStack overflowX={"auto"} key={idx}>
            <HStack spacing={"30px"} w={"full"}>
              <Avatar
                src={`https://api.dicebear.com/6.x/notionists/svg?seed=${m.user.email}`}
                borderWidth={"1px"}
                borderColor={"blackAlpha.200"}
                backgroundColor={"white"}
              />
              <Box>
                <Heading size={"sm"}>{m.user.fullName}</Heading>
                <Text>{m.user.email}</Text>
              </Box>
            </HStack>
            <Spacer />
            <HStack w="full">
              <Spacer />
              <Button variant={"outline"}>
                <Text>{m.membershipType}</Text>
              </Button>
            </HStack>
            <Spacer />
            <HStack w="full">
              <Spacer />
              <Box>
                <IconButton
                  aria-label="remove member"
                  colorScheme="red"
                  variant={"outline"}
                  icon={<UserXIcon />}
                  isDisabled={true}
                />
              </Box>
            </HStack>
          </HStack>
        ))}

        {team?.invitees.map((invite, idx) => (
          <HStack overflowX={"auto"} key={idx}>
            <HStack spacing={"30px"} w="full">
              <Avatar
                src={`https://api.dicebear.com/6.x/notionists/svg?seed=${invite.inviteeEmail}`}
                borderWidth={"1px"}
                borderColor={"blackAlpha.200"}
                backgroundColor={"white"}
              />
              <Box>
                <Text>{invite.inviteeEmail}</Text>
              </Box>
            </HStack>
            <Spacer />
            <HStack w="full">
              <Spacer />
              <Button variant={'outline'}>
                <Text>{'INVITED'}</Text>
              </Button>
            </HStack>
            <Spacer />
            <HStack w="full">
              <Spacer />
              <Box>
                <DeleteInvite teamSlug={teamSlug} inviteSlug={invite.slug} refetch={refetch} />
              </Box>
            </HStack>
          </HStack>
        ))}

      </Stack>

    </VStack>

  );
};

interface DeleteInviteProps {
  teamSlug: string;
  inviteSlug: string;
  refetch: () => void;
};

const DeleteInvite: React.FC<DeleteInviteProps> = ({ refetch, teamSlug, inviteSlug }) => {
  const [deleteInvite, { loading, error, data }] = useMutation(DELETE_INVITE);
  const deleteEmailInvite = async () => {
    const res = await deleteInvite({ variables: { teamSlug, inviteSlug } });
    if (res) refetch();
  };
  const { onOpen, isOpen, onClose } = useDisclosure();

  const toast = useToast();

  useEffect(() => {
    if (!loading && !error && data) {
      toast({
        title: 'Invite deleted successfully',
        status: 'success',
        duration: 6000,
        isClosable: true,
        position: "top",
        containerStyle: {
          paddingTop: "30px"
        },
      });
    } else if (error) {
      toast({
        title: 'Could not delete invite',
        status: 'error',
        duration: 6000,
        isClosable: true,
        position: "top",
        containerStyle: {
          paddingTop: "30px"
        },
      });
    }
  }, [loading, data, error]);

  return (
    <>
      <IconButton
        aria-label="remove invite"
        colorScheme="red"
        variant={"outline"}
        icon={<UserXIcon />}
        isDisabled={loading}
        onClick={onOpen}
      />
      <SingleActionModal
        heading={"Delete invite"}
        action={deleteEmailInvite}
        isOpen={isOpen}
        onClose={onClose}
        loading={loading}
      >
        <Text>
          Are you sure you want to delete this invite?
          You will have to invite this person again later if you want them to join your team.
        </Text>
      </SingleActionModal>
    </>
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
    case 'members':
      defaultIndex = 1;
      break;
    case 'subscription':
      defaultIndex = 2;
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
