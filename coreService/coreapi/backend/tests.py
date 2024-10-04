from django.test import TestCase
import services
from context import context as ctx

def test_video():
    mapping = services.gen_content(
        """Hey, TikTok! Ever wondered about the mighty Roman Empire? Let's dive into some quick facts!1. The Roman Empire lasted over 500 years, making it one of the most powerful civilizations in history! 
        2. Geographically, it stretched across three continents, from Britain to Egypt, creating a diverse and influential empire. 3. Julius Caesar, a key figure, played a pivotal role in the transition from the Roman Republic to the Roman Empire.
        4. The Colosseum, an iconic amphitheater, hosted gladiator battles and chariot races, entertaining thousands. 5. Romans were incredible engineers, constructing aqueducts and roads like the famous Appian Way. 6. Latin, the language of the Romans, left a lasting impact on many modern languages, including English. 7. The fall of the empire was marked by invasions, economic challenges, and political instability.
        8. Despite its fall, Rome's legacy endures, influencing our laws, architecture, and more. If you found these Roman Empire facts fascinating, hit that like button and follow for more history content! Stay curious, TikTok fam! #RomanEmpire #HistoryFacts #EducationalTikTok
        """ 
    )
    return services.gen_video(mapping)
 
def test_problem_gen():
    out = services.chat_generate_script("discrete-math", ctx)
    print(out)


if (__name__=="__main__"):
    test_problem_gen()
