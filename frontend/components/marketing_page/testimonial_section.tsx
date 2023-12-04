import { Avatar, Box, HStack, Heading, Stack, Text } from "@chakra-ui/react";
import { ExternalLink } from "lucide-react";
import Link from "next/link";

interface TestimonialCardProps {
  name: string;
  title: string;
  src: string;
  text: string[];
  link: string;
};

const TestimonialCard: React.FC<TestimonialCardProps> = (props) => {
  const { name, title, src, text, link } = props;

  return (
    <Box
      flex="1"
      p={6}
      shadow="lg"
      borderRadius="md"
      borderWidth={1}
      position="relative"
    >
      <HStack mb="30px">
        <Avatar src={src} name={name} />
        <Stack spacing={1}>
          <HStack>
            <Heading
              size={{ base: "md" }}
              fontWeight={'medium'}
            >
              { name }
            </Heading>
            <Link
              href={link}
              target='_blank'
              aria-label={`Read more about ${name}`}
            >
              <ExternalLink />
            </Link>
          </HStack>
          <Heading
            size={{ base:"sm" }}
            fontWeight={'normal'}
          >
            { title }
          </Heading>
        </Stack>
      </HStack>

      <Heading
        fontSize={"md"}
        fontWeight={"normal"}
        as={"em"}
      >
        {text.map((txt, idx) => (
          <Text pt="10px" key={idx}>{'"'}{txt}{'"'}</Text>
        ))}
      </Heading>
    </Box>
  );
};


const TestimonialSection = () => {
  return (
    <Stack
      display={"flex"}
      alignItems={{ base: "center" }}
      maxW={"1400px"}
      w={"full"}
    >
      <Heading
        fontWeight={"medium"}
        size={{ base: '2xl', md: "3xl" }}
        textAlign={"left"}
        w={"full"}
      >
        Welcome to efficient {' '}
        <Text
          as={"span"}
          bgGradient={'linear(to-tr, #007CF0, #01DFD8)'}
          bgClip='text'
        >
          broadcasting
        </Text>
      </Heading>
      <Heading
        fontWeight={'normal'}
        size={{ base: "sm", sm: "lg" }}
        textAlign={"left"}
        w={"full"}
      >
        Discover how our users are revolutionizing their content reach
      </Heading>

      <Stack
        direction={{ base: "column", md: "row" }}
        justifyContent={{ md: "space-between" }}
        w="full"
        mt="45px"
        px="16px"
        spacing={"80px"}
      >
        <TestimonialCard
          name='Rahul Pandey'
          title='Educator'
          src={'/rahulimg.jpeg'}
          text={['Amazed by this 🤯😮', 'The next few years for creators are going to be wild.']}
          link='https://www.youtube.com/@RahulInHindi'
        />
        <TestimonialCard
          name='Devin Estopinal'
          title='Social Media Content Strategist'
          src={'/devnimg.jpeg'}
          text={["I can't wait to see the results", "Just got done downloading the video after dubbing in spanish and it is absolutely flawless"]}
          link='https://x.com/NotDevn'
        />
        <TestimonialCard
          name='Dallon Asnes'
          title='Travel Content Creator'
          src={'/dallonimg.jpeg'}
          text={['first vid in progress', 'the hindi dubbing is excellent']}
          link='https://www.youtube.com/@dallonearth'
        />
      </Stack>
    </Stack>
  );
};

export default TestimonialSection;
