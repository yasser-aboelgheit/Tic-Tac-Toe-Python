import random
from string import*
G=[0,1,2,3,4,5,6,7,8]
def welcome():    
    user=input("Player Name :")
    print ("Hi %s hope you are Ready !"%(user))
    print ("Game level ?!")

#AT THIS STAGE THE USER CHOOSE THE GAME LEVEL EITHER NORMAL OR MISSION IMPOSSIBLE
    while True:
        try:
            choice = int(input("choose 1 or 2:-\n1)Normal\n2)Mission impossible\nur choice?\n"))
        except Exception:
            print("Sorry, I didn't understand that.") #the user entered non-integer value
            continue
        if choice ==1:
            Normal()
            break
        elif choice==2:
            mission_impossible()
            break
        else:
            print("You should choose either 1 or 2")#The user entered values that are not either 1 or 2
            continue

# This function is used to display the tic tac toe playfield
def show():
    print ("\u200c",G[0],"|",G[1],"|",G[2],"\n",G[3],"|",G[4],"|",G[5],"\n",G[6],"|",G[7],"|",G[8])
 
def mission_impossible():
    iteration=0
    while True:
        while True:
            show()            
            try:
                u=int(input("choice: "))
            except Exception:
                print("Invalid input, use the numbers to choose between 0 and 8") #the user entered non-integer value
                continue            
            if int(u) in range(0,9):
                u=int(u)
                break
            else:
                print ("Invalid input, choose from 0 to 8")
        
        if type(u)!=type(G[u]):
            print ("already taken Try Again \n==============")
        else:
            print ("==============")
            G[u]="X"
            if check(G)!=None:
                show()
                print (check(G))
                break
            #Computer reply
            #first check if there is a chance to win or if the user would win in the next move, so the pirority to block his chance
            if type(obstruct(G))==int:
                e=obstruct(G)
            # if the center is not taken yet, the computer choose the center
            elif G[4]==4:
                e=4
            #at the first iteration,second choose one of the corners, in case the center is already chosen
            elif G[4]=="X" and iteration<=1 :
                r=[0,2,6,8]
                e=random.sample(r,1)[0]
                while G[e]=="X" or G[e]=="O":
                    e=random.sample(r,1)[0]
            elif G[4]=="O" and iteration==1 and G[0:6].count("X")==2 and G[1:6:3].count("X")==1:
                r=[0,2]
                e=random.sample(r,1)[0]
                while G[e]=="X" or G[e]=="O":
                    e=random.sample(r,1)[0]                    
            elif G[4]=="O" and iteration==1 and G[3:9].count("X")==2 and G[1:8:3].count("X")==1:
                e=random.sample([6,8],1)[0]
                while G[e]=="X" or G[e]=="O":
                    e=random.sample(r,1)[0]
            elif G[4]=="O" and iteration==1 and G[3:9].count("X")==2:
                #special case
                e=1
            elif G[4]=="O" and iteration==1 and G[0:6].count("X")==2:
                #special case
                e=7                
            elif G[4]=="O" and iteration==1 and ((G[0:2].count("X")+G[3:5].count("X")+G[6:8].count("X"))==2):
                e=5                    
            elif G[4]=="O" and iteration==1 and G[1:3].count("X")+G[4:6].count("X")+G[7:9].count("X")==2:
                e=3                          
            elif G[4]=="O" and iteration==1:
                r=[1,3,5,7]
                e=random.sample(r,1)[0]
                while G[e]=="X":
                    e=random.sample(r,1)[0]                   
            else:
                r=range(0,9)
                e=random.sample(r,1)[0]
                while G[e]=="X" or G[e]=="O":
                    e=random.sample(r,1)[0]
                    
            G[e]="O"
            print ("computer chose %d"%(e))
        
        if check(G)!=None:
            show()
            print (check(G))
            break     
        iteration=iteration+1
        
def Normal():   
    while True:      
        while True:
            show()            
            try:
                u=int(input("choice: "))
            except Exception:
                print("Invalid input, use the numbers to choose between 0 and 8") #the user entered non-integer value
                continue            
            if int(u) in range(0,9):
                u=int(u)
                break
            else:
                print ("Invalid input, choose from 0 to 8")        
        if type(u)!=type(G[u]):
            print ("already taken Try Again \n==============")
        else:
            print ("==============")
            G[u]="X"
            if check(G)!=None:
                show()
                print (check(G))
                break
            if type(obstruct(G))==int:
                e=obstruct(G)
            else:
                r=range(0,9)
                y=random.sample(r,1)
                e=y[0]
                while G[e]=="X" or G[e]=="O":
                    y=random.sample(r,1)
                    e=y[0]
            G[e]="O"             
            print ("computer chose %d"%(e))


def check(G):
    D=[G[2:7:2],G[0:9:4],G[0:9:3],G[1:9:3],G[2:9:3],G[0:3],G[3:6],G[6:9]]
    
    c=[]
    for i in G:
        if type(i)==type(""):
            c.append(i)
    v=["X"]*3
    n=["O"]*3
    if v in D:
        return "X wins"
    elif n in D:
        return "O wins"
    elif len(G)==len(c):
        return "Draw"
    
def obstruct(G):
    D=[G[2:7:2],G[0:9:4],G[0:9:3],G[1:9:3],G[2:9:3],G[0:3],G[3:6],G[6:9]]
    
    for i in D:
        if i.count("O")==2:
            for e in i:
                if type(e)==int:
                    return e
    for i in D:
        if i.count("X")==2:
            for e in i:
                if type(e)==int: 
                    return e      
welcome()
