import {
  Box,
  Button,
  HStack,
  Heading,
  IconButton,
  VStack,
  useColorMode,
  useColorModeValue,
  useDisclosure,
} from '@chakra-ui/react';
import { useUser } from '@clerk/nextjs';
import { Moon, Sun, X, Menu } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';

const Links = [
  {
    title: 'Home',
    link: '/#',
  },
  {
    title: 'Use Cases',
    link: '/#usecases',
  },
  {
    title: 'Benefits',
    link: '/#benefits',
  },
  {
    title: 'Testimonials',
    link: '/#testimonials',
  },
  {
    title: 'Pricing',
    link: '/#pricing',
  },
];

export const NavLink = ({ children }: { children: React.ReactNode }) => (
  <Button
    px={2}
    py={1}
    rounded={'md'}
    variant={"ghost"}
    _hover={{
      backgroundColor: useColorModeValue("black", "white"),
      textColor: useColorModeValue("white", "black"),
    }}
  >
    {children}
  </Button>
);

interface NavbarProps {
  marketing?: boolean;
};

const Navbar: React.FC<NavbarProps> = ({ marketing }) => {

  const { toggleColorMode } = useColorMode();
  const { isSignedIn, isLoaded } = useUser();

  const { isOpen, onOpen, onClose } = useDisclosure();

  const blackWhite = useColorModeValue('black', 'white');
  const whiteBlack = useColorModeValue('white', 'black');

  return (
    <Box
      w="full"
      display={"flex"}
      alignItems={"center"}
      justifyContent={"center"}
      flexDir={"column"}
    >

      <Box
        display={"flex"}
        justifyContent={"space-between"}
        w="full"
        background={useColorModeValue("white", "black")}
        maxW={"1400px"}
      >

        <Box display={"flex"} alignItems={"center"} justifyContent={"center"}>
          <Image
            src={useColorModeValue('/planetcastlight.svg', '/planetcastdark.svg')}
            width={60}
            height={100}
            alt='planet cast logo'
          />
          <Heading
            fontSize={"30px"}
            display={{ base: "none", md:"flex" }}
            fontWeight={"medium"}
          >
            PlanetCast
          </Heading>
        </Box>

        <HStack>
          <HStack
            as={'nav'}
            spacing={4}
            display={{ base: 'none', lg: 'flex' }}
          >
            {Links.map((link) => (
              <Link href={link.link} key={link.title}>
                <NavLink>{link.title}</NavLink>
              </Link>
            ))}
          </HStack>

          <IconButton
            onClick={toggleColorMode}
            aria-label='color mode toggle'
            icon={useColorModeValue(<Moon />, <Sun />)}
            variant={"ghost"}
          />

          <IconButton
            onClick={isOpen ? onClose : onOpen}
            aria-label={'Open Menu'}
            icon={ isOpen ? <X /> : <Menu /> }
            display={{ base: 'inherit', lg: 'none' }}
            variant={"ghost"}
          />

          <Link
            href={'/dashboard'}
            hidden={!(marketing && isLoaded)}
          >
            <Button
              backgroundColor={blackWhite}
              textColor={whiteBlack}
              borderWidth={"1px"}
              _hover={{
                backgroundColor: whiteBlack,
                textColor: blackWhite,
              }}
            >
              { isSignedIn ? 'Dashboard' : 'Start for Free' }
            </Button>
          </Link>

        </HStack>
      </Box>

      {isOpen ? (
        <Box pb={4} display={{ lg: 'none' }} w="full">
          <VStack
            as={'nav'}
            spacing={4}
            textColor={blackWhite}
            fontWeight="600"
            alignItems={'left'}
            marginLeft={'25px'}
            marginTop={'20px'}
          >
            {Links.map((link) => (
              <Link href={link.link} key={link.title} onClick={onClose}>
                <NavLink>{link.title}</NavLink>
              </Link>
            ))}
          </VStack>
        </Box>
      ) : null}

    </Box>
  );
}

export default Navbar;
