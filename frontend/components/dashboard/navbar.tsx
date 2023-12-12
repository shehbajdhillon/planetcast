import {
  Box,
  HStack,
  MenuItem,
  Menu,
  MenuDivider,
  Avatar,
  MenuButton,
  Center,
  MenuList,
  Text,
  useColorMode,
  useColorModeValue,
  Divider,
  IconButton,
  Skeleton,
} from '@chakra-ui/react';
import { useClerk, useUser } from '@clerk/nextjs';
import { ChevronsUpDown, Clock } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { Project, Team } from '@/types';
import Button from '../button';


export const MenuBar: React.FC = () => {

  const { toggleColorMode } = useColorMode();
  const { signOut } = useClerk();
  const { isLoaded, isSignedIn } = useUser();

  const { user } = useUser();

  const router = useRouter();

  return (
    <Menu>
      <MenuButton
        hidden={!isSignedIn || !isLoaded}
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
            size={{ base: "sm" }}
            src={`https://api.dicebear.com/6.x/notionists/svg?seed=${user?.primaryEmailAddress?.emailAddress}`}
            borderRadius={"full"}
            backgroundColor={"white"}
          />
        </HStack>
      </MenuButton>
      <MenuList alignItems={'center'} backgroundColor={useColorModeValue("white", "black")}>
        <Center>
          <Avatar
            borderWidth={"1px"}
            borderColor={"blackAlpha.200"}
            size={'2xl'}
            src={`https://api.dicebear.com/6.x/notionists/svg?seed=${user?.primaryEmailAddress?.emailAddress}`}
            backgroundColor={"white"}
          />
        </Center>
        <Center maxW={"220px"} flexDir={"column"}>
          <Text
            display={{ base: "none", lg: "inline-block" }}
            whiteSpace={"nowrap"}
            textOverflow={"ellipsis"}
            overflow={'hidden'}
            noOfLines={1}
          >
            {user?.fullName}
          </Text>
          <Text
            display={{ base: "none", lg: "inline-block" }}
            whiteSpace={"nowrap"}
            textOverflow={"ellipsis"}
            overflow={'hidden'}
            noOfLines={1}
          >
            {user?.primaryEmailAddress?.emailAddress}
          </Text>
        </Center>
        <MenuDivider />
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
          onClick={() => router.push('/dashboard/account')}
          _hover={{
            backgroundColor: useColorModeValue("blackAlpha.200", "whiteAlpha.200")
          }}
        >
          Account Settings
        </MenuItem>
        <MenuItem
          backgroundColor={useColorModeValue("white", "black")}
          onClick={() => signOut({})}
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

interface NavbarProps {
  teams?: Team[];
  projects?: Project[];
  teamSlug?: string;
  projectId?: number;
};

const Navbar: React.FC<NavbarProps> = ({ teams, projects, teamSlug, projectId }) => {

  const bgColor = useColorModeValue("white", "black");
  const hoverColor = useColorModeValue("blackAlpha.200", "whiteAlpha.200");

  const router = useRouter();

  const currentTeam: Team | undefined = teams?.filter((team: Team) => team.slug === teamSlug)[0];

  return (
    <Box w="full" display={"flex"} alignItems={"center"} justifyContent={"center"}>
      <Box
        display={"flex"}
        justifyContent={"space-between"}
        w="full"
        background={useColorModeValue("white", "black")}
        maxW={"1920px"}
      >
        <HStack display={"flex"} alignItems={"center"} justifyContent={"center"} spacing={4}>

          <Link href={`/dashboard`}>
            <Image
              src={useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg')}
              width={60}
              height={100}
              alt='planet cast logo'
            />
          </Link>

          { teamSlug && <Divider orientation='vertical' borderWidth={"1px"} maxH={"40px"} transform={"rotate(20deg)"} /> }

          { teamSlug && teams ?
            <HStack display={"flex"} alignItems={"center"} justifyContent={"center"} spacing={2} h="full">
              {teams?.filter((team: Team) => team.slug === teamSlug).map((team: Team, idx: number) => (
                <Link href={`/dashboard/${team.slug}`} key={idx}>
                  <Text
                    noOfLines={1}
                    fontWeight={"medium"}
                    maxWidth={{
                      base: projectId ? "20px" : "127px",
                      sm: projectId ? "68px" : "232px",
                      md: "500px",
                    }}
                  >
                    {team.name}
                  </Text>
                </Link>
              ))}
              <Menu>
                <MenuButton as={IconButton} size={"sm"} variant={"ghost"} icon={<ChevronsUpDown />} />
                <MenuList alignItems={'center'} backgroundColor={bgColor}>
                  {teams?.map((team: Team, idx: number) => (
                    <Box key={idx}>
                      <MenuItem
                        backgroundColor={bgColor}
                        _hover={{
                          backgroundColor: hoverColor
                        }}
                        key={idx}
                        onClick={() => router.push(`/dashboard/${team.slug}`)}
                      >
                        {team.name}
                      </MenuItem>
                      { idx < teams.length - 1 && <MenuDivider /> }
                    </Box>
                  ))}
                </MenuList>
              </Menu>
            </HStack>
            :
            teamSlug && <Skeleton rounded={"lg"} w="232px" h="40px" />
          }

          { projectId && projects &&
            <Divider orientation='vertical' borderWidth={"1px"} maxH={"40px"} transform={"rotate(20deg)"} />
          }

          { projectId && projects &&
            <HStack display={"flex"} alignItems={"center"} justifyContent={"center"} spacing={4} h="full">
              {projects?.filter((project: Project) => project.id === projectId).map((project: Project, idx: number) => (
                <Link href={`/dashboard/${teamSlug}/project/${projectId}`} key={idx}>
                  <Text
                    noOfLines={1}
                    fontWeight={"medium"}
                    maxWidth={{
                      base: "20px",
                      sm: "68px",
                      md: "500px",
                    }}
                  >
                    {project.title}
                  </Text>
                </Link>
              ))}
              <Menu>
                <MenuButton as={IconButton} size={"sm"} variant={"ghost"} icon={<ChevronsUpDown />} />
                <MenuList alignItems={'center'} backgroundColor={bgColor}>
                  {projects?.map((project: Project, idx: number) => (
                    <Box key={idx}>
                      <MenuItem
                        backgroundColor={bgColor}
                        _hover={{
                          backgroundColor: hoverColor
                        }}
                        key={idx}
                        onClick={() => router.push(`/dashboard/${teamSlug}/project/${project.id}`)}
                      >
                        {project.title}
                      </MenuItem>
                      { idx < projects.length - 1 && <MenuDivider /> }
                    </Box>
                  ))}
                </MenuList>
              </Menu>
            </HStack>
          }

        </HStack>
          <HStack spacing={3}>
            { teams && teamSlug &&
              <Button
                onClick={() => router.push(`/dashboard/${teamSlug}/settings/subscription`)}
                size="sm"
              >
                <Text display={{ md: "none" }}>
                  { currentTeam?.subscriptionPlans?.[0].remainingCredits }
                </Text>
                <Box mr="5px" display={{ base: "none", md: "block" }}>
                  <Clock size={"18px"} />
                </Box>
                <Text display={{ base: "none", md: "block" }}>
                  { currentTeam?.subscriptionPlans?.[0].remainingCredits } Min
                </Text>
              </Button>
            }
            <MenuBar />
          </HStack>
      </Box>
    </Box>
  );
}

export default Navbar;
