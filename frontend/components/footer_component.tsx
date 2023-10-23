'use client'

import {
  Box,
  chakra,
  Container,
  Heading,
  Stack,
  Text,
  useColorModeValue,
  VisuallyHidden,
} from '@chakra-ui/react'
import { ReactNode } from 'react'

import Image from 'next/image';

const Logo = (props: any) => {
  return (
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
  )
}

const SocialButton = ({
  children,
  label,
  href,
}: {
  children: ReactNode
  label: string
  href: string
}) => {
  return (
    <chakra.button
      bg={useColorModeValue('blackAlpha.100', 'whiteAlpha.100')}
      rounded={'full'}
      w={8}
      h={8}
      cursor={'pointer'}
      as={'a'}
      href={href}
      display={'inline-flex'}
      alignItems={'center'}
      justifyContent={'center'}
      transition={'background 0.3s ease'}
      _hover={{
        bg: useColorModeValue('blackAlpha.200', 'whiteAlpha.200'),
      }}>
      <VisuallyHidden>{label}</VisuallyHidden>
      {children}
    </chakra.button>
  )
}

const FooterComponent = () => {
  return (
    <Box
      w="full"
      color={useColorModeValue('gray.700', 'gray.200')}>
      <Container
        as={Stack}
        maxW={'6xl'}
        py={4}
        spacing={4}
        justify={'center'}
        align={'center'}>
        <Logo />
        <Stack direction={'row'} spacing={6}>
          <Box as="a" href={'#'}>
            Home
          </Box>
          <Box as="a" href={'#benefits'}>
            Benefits
          </Box>
          {/**
          <Box as="a" href={'#testimonials'}>
            Testimonials
          </Box>
          **/}
          <Box as="a" href={'#pricing'}>
            Pricing
          </Box>
        </Stack>
      </Container>

      <Box
        borderTopWidth={1}
        borderStyle={'solid'}
      >
        <Container
          as={Stack}
          maxW={'6xl'}
          py={4}
          direction={{ base: 'column', md: 'row' }}
          spacing={4}
          justify={{ base: 'center', md: 'space-between' }}
          align={{ base: 'center', md: 'center' }}>
          <Text>© {new Date().getFullYear()} PlanetCast. All rights reserved</Text>
          <Stack direction={'row'} spacing={6}>
            <SocialButton label={'Twitter'} href={'#'}>
            </SocialButton>
            <SocialButton label={'YouTube'} href={'#'}>
            </SocialButton>
            <SocialButton label={'Instagram'} href={'#'}>
            </SocialButton>
          </Stack>
        </Container>
      </Box>
    </Box>
  )
}

export default FooterComponent;
