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
  Button,
} from '@chakra-ui/react';
import { useClerk, useUser } from '@clerk/nextjs';
import { ChevronsUpDown } from 'lucide-react';
import Image from 'next/image';
import { Dispatch, SetStateAction, useState } from 'react';

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
  selectedTeam: number;
  setSelectedTeam: Dispatch<SetStateAction<number>>;
};

const Navbar: React.FC<NavbarProps> = (props) => {

  const { teams, selectedTeam, setSelectedTeam } = props;

  const bgColor = useColorModeValue("white", "black");
  const hoverColor = useColorModeValue("blackAlpha.200", "whiteAlpha.200");

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
          <Divider orientation='vertical' borderWidth={"1px"} maxH={"40px"} transform={"rotate(20deg)"} />
          <Menu>
            <MenuButton as={Button} w="full" variant={"outline"} rightIcon={<ChevronsUpDown />}>
              {teams[selectedTeam].name}
            </MenuButton>
            <MenuList alignItems={'center'} backgroundColor={useColorModeValue("white", "black")}>
              {teams.map((team, idx) => (
                <Box key={idx}>
                  <MenuItem
                    backgroundColor={bgColor}
                    _hover={{
                      backgroundColor: hoverColor
                    }}
                    onClick={() => setSelectedTeam(idx)}
                    key={idx}
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
        <HStack>
          <MenuBar />
        </HStack>
      </Box>
    </Box>
  );
}

export default Navbar;
