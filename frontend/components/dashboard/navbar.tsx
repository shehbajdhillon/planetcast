import { gql, useQuery } from '@apollo/client';
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
import { Team } from '@/types';

const GET_TEAMS = gql`
  query GetTeams {
    getTeams {
      id
      slug
      name
      teamType
      created
    }
  }
`;

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

const Navbar: React.FC = () => {

  const bgColor = useColorModeValue("white", "black");
  const hoverColor = useColorModeValue("blackAlpha.200", "whiteAlpha.200");

  const { data } = useQuery(GET_TEAMS);


  const router = useRouter();

  const teamSlug = router.query.teamSlug;
  const projectId = router.query.projectId?.[0];

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
    console.log({ teamSlug, projectId });
  }, [teamSlug, projectId]);

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
              {data?.getTeams?.filter((team: Team) => team.slug === teamSlug).map((team: Team, idx: number) => (
                <Link href={`/dashboard/${team.slug}`} key={idx}>
                  <Text>
                    {team.name}
                  </Text>
                </Link>
              ))}
              <Menu>
                <MenuButton as={IconButton} variant={"ghost"} icon={<ChevronsUpDown />} />
                <MenuList alignItems={'center'} backgroundColor={bgColor}>
                  {data?.getTeams?.map((team: Team, idx: number) => (
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
                      <MenuDivider />
                    </Box>
                  ))}
                  <MenuItem
                    backgroundColor={bgColor}
                    _hover={{
                      backgroundColor: hoverColor
                    }}
                  >
                    Create New Team
                  </MenuItem>
                </MenuList>
              </Menu>
            </HStack>
          }

          { projectId &&
            <HStack display={"flex"} alignItems={"center"} justifyContent={"center"} spacing={4} h="full">
              <Divider orientation='vertical' borderWidth={"1px"} maxH={"40px"} transform={"rotate(20deg)"} />
              {data?.getTeams?.filter((team: Team) => team.slug === teamSlug).map((team: Team, idx: number) => (
                <Link href={`/dashboard/${teamSlug}/${team.slug}`} key={idx}>
                  <Text>
                    {team.name}
                  </Text>
                </Link>
              ))}
              <Menu>
                <MenuButton as={IconButton} variant={"ghost"} icon={<ChevronsUpDown />} />
                <MenuList alignItems={'center'} backgroundColor={bgColor}>
                  {data?.getTeams?.map((team: Team, idx: number) => (
                    <Box key={idx}>
                      <MenuItem
                        backgroundColor={bgColor}
                        _hover={{
                          backgroundColor: hoverColor
                        }}
                        key={idx}
                        onClick={() => router.push(`/dashboard/${teamSlug}/${team.slug}`)}
                      >
                        {team.name}
                      </MenuItem>
                      <MenuDivider />
                    </Box>
                  ))}
                  <MenuItem
                    backgroundColor={bgColor}
                    _hover={{
                      backgroundColor: hoverColor
                    }}
                  >
                    Create New Team
                  </MenuItem>
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
