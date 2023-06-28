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
} from '@chakra-ui/react';
import { useClerk, useUser } from '@clerk/nextjs';
import { ChevronsUpDown } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

import NProgress from 'nprogress';
import { Project, Team } from '@/types';

export const MenuBar: React.FC = () => {

  const { toggleColorMode } = useColorMode();
  const { signOut } = useClerk();

  const { user } = useUser();

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
        <Center maxW={"220px"}>
          <Text
            display={{ base: "none", lg: "inline-block" }}
            whiteSpace={"nowrap"}
            textOverflow={"ellipsis"}
            overflow={'hidden'}
            noOfLines={1}
          >
            {user?.fullName}
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
  teams: Team[];
  projects: Project[];
  teamSlug: string;
  projectId?: number;
};

const Navbar: React.FC<NavbarProps> = ({ teams, projects, teamSlug, projectId }) => {

  const bgColor = useColorModeValue("white", "black");
  const hoverColor = useColorModeValue("blackAlpha.200", "whiteAlpha.200");

  const router = useRouter();

  useEffect(() => {
    const handleRouteStart = () => NProgress.start();
    const handleRouteDone = () => NProgress.done();

    router.events.on('routeChangeStart', handleRouteStart);
    router.events.on('routeChangeComplete', handleRouteDone);
    router.events.on('routeChangeError', handleRouteDone);

    return () => {
      handleRouteDone();
      router.events.off('routeChangeStart', handleRouteStart);
      router.events.off('routeChangeComplete', handleRouteDone);
      router.events.off('routeChangeError', handleRouteDone);
    };
  }, [router.events]);

  useEffect(() => {
    console.log({ teamSlug, projectId, teams, projects });
  }, [teamSlug, projectId, teams, projects]);

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
          <Image
            src={useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg')}
            width={60}
            height={100}
            alt='planet cast logo'
          />

          { teamSlug &&
            <HStack display={"flex"} alignItems={"center"} justifyContent={"center"} spacing={4} h="full">
              <Divider orientation='vertical' borderWidth={"1px"} maxH={"40px"} transform={"rotate(20deg)"} />
              {teams?.filter((team: Team) => team.slug === teamSlug).map((team: Team, idx: number) => (
                <Link href={`/${team.slug}`} key={idx}>
                  <Text
                    noOfLines={1}
                    maxWidth={{
                      base: projectId ? "33px" : "157px",
                      sm: projectId ? "116px" : "232px",
                      md: "500px",
                    }}
                  >
                    {team.name}
                  </Text>
                </Link>
              ))}
              <Menu>
                <MenuButton as={IconButton} variant={"ghost"} icon={<ChevronsUpDown />} />
                <MenuList alignItems={'center'} backgroundColor={bgColor}>
                  {teams?.map((team: Team, idx: number) => (
                    <Box key={idx}>
                      <MenuItem
                        backgroundColor={bgColor}
                        _hover={{
                          backgroundColor: hoverColor
                        }}
                        key={idx}
                        onClick={() => router.push(`/${team.slug}`)}
                      >
                        {team.name}
                      </MenuItem>
                      { idx < teams.length - 1 && <MenuDivider /> }
                    </Box>
                  ))}
                </MenuList>
              </Menu>
            </HStack>
          }

          { projectId &&
            <HStack display={"flex"} alignItems={"center"} justifyContent={"center"} spacing={4} h="full">
              <Divider orientation='vertical' borderWidth={"1px"} maxH={"40px"} transform={"rotate(20deg)"} />
              {projects?.filter((project: Project) => project.id === projectId).map((project: Project, idx: number) => (
                <Link href={`/${teamSlug}/${projectId}`} key={idx}>
                  <Text
                    noOfLines={1}
                    maxWidth={{
                      base: "33px",
                      sm: "116px",
                      md: "500px",
                    }}
                  >
                    {project.title}
                  </Text>
                </Link>
              ))}
              <Menu>
                <MenuButton as={IconButton} variant={"ghost"} icon={<ChevronsUpDown />} />
                <MenuList alignItems={'center'} backgroundColor={bgColor}>
                  {projects?.map((project: Project, idx: number) => (
                    <Box key={idx}>
                      <MenuItem
                        backgroundColor={bgColor}
                        _hover={{
                          backgroundColor: hoverColor
                        }}
                        key={idx}
                        onClick={() => router.push(`/${teamSlug}/${project.id}`)}
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
        <HStack>
          <MenuBar />
        </HStack>
      </Box>
    </Box>
  );
}

export default Navbar;
