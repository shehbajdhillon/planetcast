import SingleActionModal from "@/components/single_action_modal";
import { Team, TeamInvite } from "@/types";
import { gql, useMutation } from "@apollo/client";
import {
  HStack,
  Button,
  Text,
  Heading,
  Stack,
  VStack,
  Box,
  Spacer,
  Grid,
  GridItem,
  Drawer,
  DrawerOverlay,
  DrawerContent,
  DrawerCloseButton,
  DrawerHeader,
  DrawerBody,
  useDisclosure,
  IconButton,
  useToast,
} from "@chakra-ui/react";
import { ExternalLink, Menu } from "lucide-react";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";

interface TabButtonsProps {
  tabIdx: number;
  switchTab: (tabIdx: number) => void;
};

const DELETE_INVITE = gql`
  mutation DeleteInvite($inviteSlug: String!) {
    deleteTeamInvite(inviteSlug: $inviteSlug)
  }
`

const ACCEPT_INVITE = gql`
  mutation AcceptInvite($inviteSlug: String!) {
    acceptTeamInvite(inviteSlug: $inviteSlug)
  }
`

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
  refetch: () => void;
}

const TeamInvitesTab: React.FC<TeamInvitesTabProps> = (props) => {
  const { refetch, teams, invites, drawerOpen } = props;
  const router = useRouter();


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
          <Box w="full">
            <Stack spacing={"25px"}>
            <Text>Current Teams</Text>
            {teams?.map((team, idx) => (
              <HStack key={idx} w="full">
                <HStack w="full">
                  <Heading size="sm">{team.teamName}</Heading>
                </HStack>
                <Spacer />
                <HStack w="full">
                  <Button variant={"outline"} pointerEvents={"none"}>
                    <Text>{team.membershipType}</Text>
                  </Button>
                </HStack>
                <Spacer />
                <HStack>
                  <Spacer />
                  <IconButton
                    variant={"outline"}
                    aria-label="go to team"
                    icon={<ExternalLink />}
                    onClick={() => router.push(`/dashboard/${team.teamSlug}`)}
                  />
                </HStack>
              </HStack>
            ))}
            </Stack>
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
          <Box w="full">
            <Stack spacing={"25px"}>
            <Text>Current Invites</Text>
            {invites?.map((invite, idx) => (
              <HStack key={idx} w="full">
                <HStack w="full">
                  <Heading size="sm">{invite.teamName}</Heading>
                </HStack>
                <Spacer />
                <HStack w="full">
                  <Spacer />
                  <AcceptInvite
                    refetch={refetch}
                    inviteSlug={invite.inviteSlug}
                    teamSlug={invite.teamSlug}
                  />
                  <DeleteInvite
                    refetch={refetch}
                    inviteSlug={invite.inviteSlug}
                    teamSlug={invite.teamSlug}
                  />
                </HStack>
              </HStack>
            ))}
            {invites?.length === 0 && <Heading size="sm">No Active Invitations</Heading>}
            </Stack>
          </Box>
          <Spacer />
          <Box>

          </Box>
        </HStack>
      </Stack>
    </VStack>
  );
};

interface InviteActionProps {
  teamSlug: string;
  inviteSlug: string;
  refetch: () => void;
};

const DeleteInvite: React.FC<InviteActionProps> = ({ refetch, inviteSlug }) => {
  const { onOpen, isOpen, onClose } = useDisclosure();

  const toast = useToast();

  const [deleteInvite, { loading, error, data }] = useMutation(DELETE_INVITE);

  const deleteEmailInvite = async () => {
    const res = await deleteInvite({ variables: { inviteSlug } });
    if (res) refetch();
  };

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
      <Button
        onClick={onOpen}
        variant={"outline"}
        isDisabled={loading}
      >
        Delete
      </Button>
      <SingleActionModal
        heading={"Delete invite"}
        action={deleteEmailInvite}
        isOpen={isOpen}
        onClose={onClose}
        loading={loading}
      >
        <Text>
          Are you sure you want to delete this invite?
          You will need to be invited again later if you want to be a part of this team.
        </Text>
      </SingleActionModal>
    </>
  );
};

const AcceptInvite: React.FC<InviteActionProps> = ({ refetch, inviteSlug }) => {
  const { onOpen, isOpen, onClose } = useDisclosure();

  const toast = useToast();

  const [acceptInvite, { loading, data, error }] = useMutation(ACCEPT_INVITE);
  const acceptEmailInvite = async () => {
    const res = await acceptInvite({ variables: { inviteSlug } });
    if (res) refetch();
  };

  useEffect(() => {
    if (!loading && !error && data) {
      toast({
        title: 'Invite accepted successfully',
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
        title: 'Could not accept invite',
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
      <Button
        onClick={onOpen}
        variant={"outline"}
        isDisabled={loading}
      >
        Accept
      </Button>
      <SingleActionModal
        heading={"Delete invite"}
        action={acceptEmailInvite}
        isOpen={isOpen}
        onClose={onClose}
        loading={loading}
      >
        <Text>
          Are you sure you want to accept this invite?
        </Text>
      </SingleActionModal>
    </>
  );
};

interface AccountSettingsTab {
  loading: boolean;
  teams: Team[];
  invites: TeamInvite[];
  refetch: () => void;
};

const AccountSettingsTab: React.FC<AccountSettingsTab> = (props) => {

  const { refetch, teams, invites, loading } = props;
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

            {tabIdx === 0 && <TeamInvitesTab drawerOpen={onOpen} teams={teams} invites={invites} refetch={refetch} />}
          </GridItem>
        </Grid>
      </Box>
    </Box>
  );
};

export default AccountSettingsTab;
